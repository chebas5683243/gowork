package template

import (
	"fmt"
	"io"
	"local/utils"
	"os"
	"path/filepath"
	"strings"
)

const VIEW_RESOURCES_FOLDER = "resources/views"

const TEMPLATE_EXTENSION = ".goplate.html"

const EXPRESION_DELIMITER_START_CHAR byte = '{'

const EXPRESION_DELIMITER_END_CHAR byte = '}'

const REQUIRED_CONSECUTIVE_CHARS = 3

type TemplateData = map[string]interface{}

type Template struct {
	templatePath string
	data         TemplateData
}

func New(templatePath string, values TemplateData) *Template {
	return &Template{
		templatePath: templatePath,
		data:         values,
	}
}

func (t *Template) Get() (string, error) {
	dat, err := os.ReadFile(t.templatePath)

	if err != nil {
		return "", err
	}

	return string(dat), nil
}

func (t *Template) Parse() (string, error) {
	f, err := openFile(t.templatePath)

	if err != nil {
		return "", err
	}

	defer f.Close()

	fileBuffer := make([]byte, 1)
	contentBuffer := make([]byte, 1)
	bracketsBuffer := make([]byte, 0, 3)
	expressionBuffer := make([]byte, 0)

	readingExpression := false

	for {
		_, err := f.Read(fileBuffer)

		if err == io.EOF {
			break
		}

		if err != nil {
			return "", fmt.Errorf("error reading file: %v", err)
		}

		bracketsBuffer = append(bracketsBuffer, fileBuffer[0])

		// if reading expression
		if readingExpression {
			// if char is delimiter
			if fileBuffer[0] == EXPRESION_DELIMITER_END_CHAR {
				// check if the end delimiter is completed
				if len(bracketsBuffer) == REQUIRED_CONSECUTIVE_CHARS {
					// start writing to content buffer
					readingExpression = false

					// map value of the expresion
					key := strings.TrimSpace(string(expressionBuffer))
					mappedValue, exists := t.data[key]
					if !exists {
						return "", fmt.Errorf("key %s not found in template data", key)
					}

					contentBuffer = append(contentBuffer, toBytes(mappedValue)...)
					bracketsBuffer = nil
					expressionBuffer = nil
				}
				continue
			}

			// if it is no delimiter empty the buffer and add to the expresion
			expressionBuffer = append(expressionBuffer, bracketsBuffer...)
			bracketsBuffer = nil
			continue
		}

		// if reading content
		// if char is delimiter
		if fileBuffer[0] == EXPRESION_DELIMITER_START_CHAR {
			// check if the start delimiter is completed
			if len(bracketsBuffer) == REQUIRED_CONSECUTIVE_CHARS {
				// Start writing to expresion buffer
				readingExpression = true
				bracketsBuffer = nil
			}
			continue
		}

		// if it is no delimiter empty the buffer and add to the content
		contentBuffer = append(contentBuffer, bracketsBuffer...)
		bracketsBuffer = nil
	}

	return string(contentBuffer), nil
}

func toTemplatePath(templateName string) string {
	return filepath.FromSlash(strings.ReplaceAll(templateName, ".", "/"))
}

func getTemplatePath(templateName string) (string, error) {
	basePath, err := utils.GetProjectRoot()

	if (err) != nil {
		return "", fmt.Errorf("couldn't find views resources folder")
	}

	fileName := toTemplatePath(templateName) + TEMPLATE_EXTENSION
	templatePath := filepath.Join(basePath, VIEW_RESOURCES_FOLDER, fileName)

	if _, err := os.Stat(templatePath); err == nil {
		return templatePath, nil
	}

	return "", fmt.Errorf("couln't find template %v", templateName)
}

func toBytes(value interface{}) []byte {
	return []byte(fmt.Sprintf("%v", value))
}

func openFile(templateName string) (*os.File, error) {

	path, err := getTemplatePath(templateName)

	if err != nil {
		return nil, err
	}

	f, err := os.Open(path)

	if err != nil {
		return nil, fmt.Errorf("error opening template: %v", err)
	}

	return f, nil
}
