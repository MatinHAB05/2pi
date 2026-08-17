package static

import (
	"bytes"
)

func (sf *StaticFiles) ExecuteWakeUpTemplate() (*bytes.Buffer, error) { //TODO
	var htmlBuffer bytes.Buffer
	if err := sf.wakeupHTMLTemplate.Execute(&htmlBuffer, map[string]any{}); err != nil {
		return nil, err
	}
	return &htmlBuffer, nil
}
