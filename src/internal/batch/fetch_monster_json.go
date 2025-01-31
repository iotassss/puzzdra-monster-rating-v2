package batch

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

func (batch *Batch) FetchMonsterJSON() error {
	resp, err := http.Get(batch.monsterSourceDataJsonURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to download file: %s", resp.Status)
	}

	out, err := os.Create(batch.monsterDataJsonFilePath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}
