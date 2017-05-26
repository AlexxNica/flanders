package mysql

import (
	"fmt"

	"github.com/weave-lab/flanders/db"
)

func (m *MySQL) GetSettings(t string) (db.SettingResult, error) {

	rows, err := m.db.Query("SELECT `key`, `value` FROM settings WHERE type=?", t)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	results := db.SettingResult{}
	for rows.Next() {
		var r db.SettingObject
		err = rows.Scan(&r.Key, &r.Val)
		if err != nil {
			return nil, err
		}

		results = append(results, r)
	}

	return results, nil
}

// SetSetting creates or updates the setting type/key
func (m *MySQL) SetSetting(t string, s db.SettingObject) error {

	// if there is an error, we won't do anything with it
	//_ = m.DeleteSetting(t, s.Key)

	_, err := m.db.Exec(`INSERT INTO settings 
		(settings.type, settings.key, settings.value) 
		VALUES (?, ?, ?) 
		ON DUPLICATE KEY UPDATE value=?`, t, s.Key, s.Val, s.Val)
	if err != nil {
		fmt.Print(err)
		return err
	}

	return nil
}

func (m *MySQL) DeleteSetting(t string, k string) error {

	_, err := m.db.Exec("DELETE FROM settings WHERE `type`=? and `key`=?", t, k)
	if err != nil {
		return err
	}

	return nil
}
