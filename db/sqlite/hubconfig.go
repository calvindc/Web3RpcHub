package sqlite

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/calvindc/Web3RpcHub/db"
	"github.com/calvindc/Web3RpcHub/db/sqlite/models"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

var _ db.HubConfig = (*Config)(nil)

// the database will only ever store one row, which contains all the hub settings
const configRowID = 0

/* Config basically enables long-term memory for the server when it comes to storing settings. Currently, the only
* stored settings is the privacy mode of the hub.
 */
type Config struct {
	db *sql.DB
}

func (c Config) GetPrivacyMode(ctx context.Context) (db.PrivacyMode, error) {
	config, err := models.FindConfig(ctx, c.db, configRowID)
	if err != nil {
		return db.ModeUnknown, err
	}

	return config.PrivacyMode, nil
}

func (c Config) SetPrivacyMode(ctx context.Context, pm db.PrivacyMode) error {
	// make sure the privacy mode is an ok value
	err := pm.IsValid()
	if err != nil {
		return err
	}

	err = transact(c.db, func(tx *sql.Tx) error {
		// get the settings row
		config, err := models.FindConfig(ctx, tx, configRowID)
		if err != nil {
			return err
		}

		// set the new privacy mode
		config.PrivacyMode = pm
		// issue update stmt
		rowsAffected, err := config.Update(ctx, tx, boil.Infer())
		if err != nil {
			return err
		}
		if rowsAffected == 0 {
			return fmt.Errorf("setting privacy mode should have update the settings row, instead 0 rows were updated")
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil // alles gut!!
}

func (c Config) GetDefaultLanguage(ctx context.Context) (string, error) {
	config, err := models.FindConfig(ctx, c.db, configRowID)
	if err != nil {
		return "", err
	}

	return config.DefaultLanguage, nil
}

func (c Config) SetDefaultLanguage(ctx context.Context, langTag string) error {
	if len(langTag) == 0 {
		return fmt.Errorf("language tag cannot be empty")
	}

	err := transact(c.db, func(tx *sql.Tx) error {
		// get the settings row
		config, err := models.FindConfig(ctx, tx, configRowID)
		if err != nil {
			return err
		}

		// set the new language tag
		config.DefaultLanguage = langTag
		// issue update stmt
		rowsAffected, err := config.Update(ctx, tx, boil.Infer())
		if err != nil {
			return err
		}
		if rowsAffected == 0 {
			return fmt.Errorf("setting default language should have update the settings row, instead 0 rows were updated")
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}
