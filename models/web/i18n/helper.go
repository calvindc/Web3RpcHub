package i18n

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"io/fs"

	"path/filepath"

	"sort"

	"net/http"

	"html/template"

	"github.com/BurntSushi/toml"
	"github.com/calvindc/Web3RpcHub/db"
	"github.com/calvindc/Web3RpcHub/internal/repository"
	"github.com/calvindc/Web3RpcHub/models/web/utils"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"go.mindeco.de/http/render"
	"golang.org/x/text/language"
)

// CookieMaxAge unit:second
const CookieMaxAge = 1 * 60 * 60

type LangHelper struct {
	bundle      *i18n.Bundle
	languages   []TagTranslation
	CookieCodec []securecookie.Codec
	CookieStore *sessions.CookieStore
	config      db.HubConfig //系统的设置
}

// Localizer
type Localizer struct {
	loc *i18n.Localizer
}

// TagTranslation语种和对应的译文
type TagTranslation struct {
	Tag         string
	Translation string
}

func New(r repository.Interface, config db.HubConfig) (*LangHelper, error) {
	//构建server cookie
	cookieCodec, err := utils.LoadOrCreateCookieSecrets(r)
	if err != nil {
		return nil, err
	}
	cookieStore := &sessions.CookieStore{
		Codecs: cookieCodec,
		Options: &sessions.Options{
			Path:   "/",
			MaxAge: CookieMaxAge, //有效期一个小时
		},
	}

	//构建语言环境
	bundle := i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)

	// parse toml files and add them to the bundle
	walkFn := func(path string, info os.FileInfo, rs io.Reader, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		if !strings.HasSuffix(path, "toml") {
			return nil
		}

		mfb, err := ioutil.ReadAll(rs)
		if err != nil {
			return err
		}
		_, err = bundle.ParseMessageFileBytes(mfb, path)
		if err != nil {
			return fmt.Errorf("i18n: failed to parse file %s: %w", path, err)
		}
		fmt.Println("i18n loaded language toml file ", path)
		return nil
	}
	// walk the embedded defaults
	err = fs.WalkDir(Defaults, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		info, err := d.Info()
		if err != nil {
			return err
		}

		r, err := Defaults.Open(path)
		if err != nil {
			return err
		}

		err = walkFn(path, info, r, err)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("i18n: failed to iterate localizations: %w", err)
	}

	// walk the local filesystem for overrides and additions
	err = filepath.Walk(r.GetPath("i18n"), func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		r, err := os.Open(path)
		if err != nil {
			return err
		}
		defer r.Close()

		err = walkFn(path, info, r, err)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil && !os.IsNotExist(err) {
		return nil, fmt.Errorf("i18n: failed to iterate localizations: %w", err)
	}

	langmap := listLanguages(bundle)

	return &LangHelper{
		bundle:      bundle,
		languages:   langmap,
		CookieCodec: cookieCodec,
		CookieStore: cookieStore,
		config:      config,
	}, nil
}

// listLanguages 构建语言bundle content :按语言tag字母顺序
func listLanguages(bundle *i18n.Bundle) []TagTranslation {
	languageTags := bundle.LanguageTags()
	tags := make([]string, 0, len(languageTags))

	languageslice := make([]TagTranslation, 0, len(languageTags))

	for _, langTag := range languageTags {
		tags = append(tags, langTag.String())
	}
	// 按首字母排序
	sort.Strings(tags)

	// 按顺序把语言标签和翻译进行绑定到 []TagTranslation  0:chinese-你好 1:en-hello
	for _, langTag := range tags {
		var l Localizer
		l.loc = i18n.NewLocalizer(bundle, langTag)

		msg, err := l.loc.Localize(&i18n.LocalizeConfig{
			MessageID: "LanguageName",
		})
		if err != nil {
			msg = langTag
		}

		languageslice = append(languageslice, TagTranslation{Tag: langTag, Translation: msg})
	}

	return languageslice
}

func (h LangHelper) ListLanguages() []TagTranslation {
	return h.languages
}

func (h LangHelper) ChooseTranslation(requestedTag string) string {
	for _, xentry := range h.languages {
		if xentry.Tag == requestedTag {
			return xentry.Translation
		}
	}
	return requestedTag
}

func (h LangHelper) newLocalizer(lang string, accept ...string) *Localizer {
	var langs = []string{lang}
	langs = append(langs, accept...)
	var l Localizer
	l.loc = i18n.NewLocalizer(h.bundle, langs...)
	return &l
}

const LanguageCookieName = "hub-language"

// FromRequest 通过调用方的lang定义翻译资源
func (h LangHelper) FromRequest(r *http.Request) *Localizer {
	lang := r.FormValue("lang")
	accept := r.Header.Get("Accept-Language")

	session, err := h.CookieStore.Get(r, LanguageCookieName)
	if err != nil {
		return h.newLocalizer(lang, accept)
	}

	prevCookie := session.Values["lang"]
	if prevCookie != nil {
		return h.newLocalizer(prevCookie.(string), lang, accept)
	}

	defaultLang, err := h.config.GetDefaultLanguage(r.Context())
	if err != nil {
		return h.newLocalizer(lang, accept)
	}
	// 如果没有获取到用户的cookie设置, 则使用hub的默认设置
	return h.newLocalizer(defaultLang, accept)
}

// -----------------------------------------------------------------------------------

func (h LangHelper) GetRenderFuncs() []render.Option {
	var opts = []render.Option{
		render.InjectTemplateFunc("i18npl", func(r *http.Request) interface{} {
			loc := h.FromRequest(r)
			return loc.LocalizePlurals
		}),

		render.InjectTemplateFunc("i18n", func(r *http.Request) interface{} {
			loc := h.FromRequest(r)
			return loc.LocalizeSimple
		}),

		render.InjectTemplateFunc("i18nWithData", func(r *http.Request) interface{} {
			loc := h.FromRequest(r)
			return loc.LocalizeWithData
		}),
	}
	return opts
}

func (l Localizer) LocalizeSimple(messageID string) template.HTML {
	msg, err := l.loc.Localize(&i18n.LocalizeConfig{
		MessageID: messageID,
	})
	if err == nil {
		return template.HTML(msg)
	}

	panic(fmt.Sprintf("i18n/error: failed to localize label %s: %s", messageID, err))
}

func (l Localizer) LocalizeWithData(messageID string, labelsAndData ...string) template.HTML {
	n := len(labelsAndData)
	if n%2 != 0 {
		panic(fmt.Errorf("expected an even amount of labels and data. got %d", n))
	}

	tplData := make(map[string]string, n/2)
	for i := 0; i < n; i += 2 {
		key := labelsAndData[i]
		data := labelsAndData[i+1]
		tplData[key] = data
	}

	msg, err := l.loc.Localize(&i18n.LocalizeConfig{
		MessageID:    messageID,
		TemplateData: tplData,
	})
	if err == nil {
		return template.HTML(msg)
	}

	panic(fmt.Sprintf("i18n/error: failed to localize label %s: %s", messageID, err))
}

func (l Localizer) LocalizePlurals(messageID string, pluralCount int) template.HTML {
	msg, err := l.loc.Localize(&i18n.LocalizeConfig{
		MessageID:   messageID,
		PluralCount: pluralCount,
		TemplateData: map[string]int{
			"Count": pluralCount,
		},
	})
	if err == nil {
		return template.HTML(msg)
	}

	panic(fmt.Sprintf("i18n/error: failed to localize label %s: %s", messageID, err))
}

func (l Localizer) LocalizePluralsWithData(messageID string, pluralCount int, tplData map[string]string) template.HTML {
	msg, err := l.loc.Localize(&i18n.LocalizeConfig{
		MessageID:    messageID,
		PluralCount:  pluralCount,
		TemplateData: tplData,
	})
	if err == nil {
		return template.HTML(msg)
	}

	panic(fmt.Sprintf("i18n/error: failed to localize label %s: %s", messageID, err))
}
