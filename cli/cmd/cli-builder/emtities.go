package main

import "github.com/alioth-center/infrastructure/cli"

// template variables for handlers/init.gen.go
type (
	HandlersInitModel struct {
		Handlers []HandlersInitItemModel
	}

	HandlersInitItemModel struct {
		Name string
	}
)

// template variables for injectors/init.gen.go
type (
	InjectorsInitModel struct {
		Injectors []InjectorsInitItemModel
	}

	InjectorsInitItemModel struct {
		Name string
	}
)

// template variables for init.gen.go and init command
type (
	InitModel struct {
		PackageName    string
		ConfigFileName string
	}

	InitCommandModel struct {
		GoBinaryPath     string
		PackageName      string
		WorkingDirectory string
	}
)

// template variables for new handlers/injectors
type (
	NewHandlerModel struct {
		HandlerName string
	}

	NewInjectorModel struct {
		InjectorName string
	}
)

type CliBuilderConfig struct {
	CliConfigFileName string            `yaml:"cli_config_file_name"`
	I18nFileName      string            `yaml:"i18n_file_name"`
	TemplateFileNames TemplateFileNames `yaml:"template_file_names"`
	OutputFiles       OutputFiles       `yaml:"output_files"`
	ProjectDirectory  string            `yaml:"project_directory"`
	PackageName       string            `yaml:"package_name"`
}

type TemplateFileNames struct {
	MainTemplate          string `yaml:"main_template"`
	InitTemplate          string `yaml:"init_template"`
	InitHandlersTemplate  string `yaml:"init_handlers_template"`
	InitInjectorsTemplate string `yaml:"init_injectors_template"`
	NewHandlerTemplate    string `yaml:"new_handler_template"`
	NewInjectorTemplate   string `yaml:"new_injector_template"`
}

type OutputFiles struct {
	MainFile          string `yaml:"main_file"`
	InitFile          string `yaml:"init_file"`
	InitHandlersFile  string `yaml:"init_handlers_file"`
	InitInjectorsFile string `yaml:"init_injectors_file"`
	NewHandlerFile    string `yaml:"new_handler_file"`
	NewInjectorFile   string `yaml:"new_injector_file"`
}

type Localization map[string][]cli.TranslatedItem

type Commands struct {
	InitProjects    string
	GetDependencies string
	FormatCode      string
}

const (
	i18nPrintWorkingDirectory                  = "print_working_directory"
	i18nInitProjectBadParams                   = "init_project_bad_params"
	i18nInitProjectReadWorkingDirectoryError   = "init_project_read_working_directory_error"
	i18nInitProjectCreateWorkingDirectoryCheck = "init_project_create_working_directory_check"
	i18nInitProjectCreateWorkingDirectoryAbort = "init_project_create_working_directory_abort"
	i18nInitProjectCreateWorkingDirectoryError = "init_project_create_working_directory_error"
	i18nInitProjectNotDirectory                = "init_project_not_directory"
	i18nInitExistProjectNotice                 = "init_exist_project_notice"
	i18nInitProjectCannotWrite                 = "init_project_cannot_write"
)
