package cli

import (
	"github.com/alioth-center/infrastructure/utils/values"
	"github.com/joeycumines/go-prompt"
)

const (
	i18nBadCommand                           = "bad_command"
	i18nNoCommand                            = "no_command"
	i18nNoImplement                          = "no_implement"
	i18nErrHandlerNotFound                   = "err_handler_not_found"
	i18nErrInjectorNotFound                  = "err_injector_not_found"
	i18nErrPreferredTranslationNotConfigured = "err_preferred_translation_not_configured"
)

var i18nPacks map[string]*TranslationSet

var fallbacks = map[string][]TranslatedItem{
	i18nBadCommand: {
		{
			Language:           "en-US",
			DisplayKey:         "Bad Command",
			DisplayDescription: "[${command}] is not a valid command",
		},
		{
			Language:           "zh-CN",
			DisplayKey:         "错误命令",
			DisplayDescription: "[${command}] 不是一个有效的命令",
		},
		{
			Language:           "ja-JP",
			DisplayKey:         "悪いコマンド",
			DisplayDescription: "[${command}] は有効なコマンドではありません",
		},
		{
			Language:           "ko-KR",
			DisplayKey:         "잘못된 명령",
			DisplayDescription: "[${command}] 는 유효한 명령이 아닙니다",
		},
		{
			Language:           "fr-FR",
			DisplayKey:         "Mauvaise Commande",
			DisplayDescription: "[${command}] n'est pas une commande valide",
		},
		{
			Language:           "es-ES",
			DisplayKey:         "Comando Incorrecto",
			DisplayDescription: "[${command}] no es un comando válido",
		},
		{
			Language:           "pt-BR",
			DisplayKey:         "Comando Ruim",
			DisplayDescription: "[${command}] não é um comando válido",
		},
		{
			Language:           "ru-RU",
			DisplayKey:         "Плохая Команда",
			DisplayDescription: "[${command}] не является допустимой командой",
		},
		{
			Language:           "de-DE",
			DisplayKey:         "Schlechter Befehl",
			DisplayDescription: "[${command}] ist kein gültiger Befehl",
		},
		{
			Language:           "it-IT",
			DisplayKey:         "Comando Errato",
			DisplayDescription: "[${command}] non è un comando valido",
		},
		{
			Language:           "nl-NL",
			DisplayKey:         "Slechte Opdracht",
			DisplayDescription: "[${command}] is geen geldige opdracht",
		},
		{
			Language:           "pl-PL",
			DisplayKey:         "Zła Komenda",
			DisplayDescription: "[${command}] nie jest prawidłową komendą",
		},
		{
			Language:           "tr-TR",
			DisplayKey:         "Kötü Komut",
			DisplayDescription: "[${command}] geçerli bir komut değil",
		},
		{
			Language:           "ar-SA",
			DisplayKey:         "أمر سيء",
			DisplayDescription: "[${command}] ليس أمرًا صالحًا",
		},
	},
	i18nNoCommand: {
		{
			Language:           "en-US",
			DisplayKey:         "No Command",
			DisplayDescription: "No such command",
		},
		{
			Language:           "zh-CN",
			DisplayKey:         "无命令",
			DisplayDescription: "没有这样的命令",
		},
		{
			Language:           "ja-JP",
			DisplayKey:         "コマンドなし",
			DisplayDescription: "そのようなコマンドはありません",
		},
		{
			Language:           "ko-KR",
			DisplayKey:         "명령 없음",
			DisplayDescription: "그런 명령이 없습니다",
		},
		{
			Language:           "fr-FR",
			DisplayKey:         "Pas de Commande",
			DisplayDescription: "Pas de telle commande",
		},
		{
			Language:           "es-ES",
			DisplayKey:         "Sin Comando",
			DisplayDescription: "No hay tal comando",
		},
		{
			Language:           "pt-BR",
			DisplayKey:         "Sem Comando",
			DisplayDescription: "Nenhum comando desse tipo",
		},
		{
			Language:           "ru-RU",
			DisplayKey:         "Нет Команды",
			DisplayDescription: "Нет такой команды",
		},
		{
			Language:           "de-DE",
			DisplayKey:         "Kein Befehl",
			DisplayDescription: "Kein solcher Befehl",
		},
		{
			Language:           "it-IT",
			DisplayKey:         "Nessun Comando",
			DisplayDescription: "Nessun comando del genere",
		},
		{
			Language:           "nl-NL",
			DisplayKey:         "Geen Opdracht",
			DisplayDescription: "Geen dergelijk commando",
		},
		{
			Language:           "pl-PL",
			DisplayKey:         "Brak Komendy",
			DisplayDescription: "Brak takiej komendy",
		},
		{
			Language:           "tr-TR",
			DisplayKey:         "Komut Yok",
			DisplayDescription: "Böyle bir komut yok",
		},
		{
			Language:           "ar-SA",
			DisplayKey:         "لا يوجد أمر",
			DisplayDescription: "لا يوجد مثل هذا الأمر",
		},
	},
	i18nNoImplement: {
		{
			Language:           "en-US",
			DisplayKey:         "No Implement",
			DisplayDescription: "This function is not implemented",
		},
		{
			Language:           "zh-CN",
			DisplayKey:         "未实现",
			DisplayDescription: "此功能未实现",
		},
		{
			Language:           "ja-JP",
			DisplayKey:         "未実装",
			DisplayDescription: "この機能は実装されていません",
		},
		{
			Language:           "ko-KR",
			DisplayKey:         "구현되지 않음",
			DisplayDescription: "이 기능은 구현되지 않았습니다",
		},
		{
			Language:           "fr-FR",
			DisplayKey:         "Non Implémenté",
			DisplayDescription: "Cette fonctionnalité n'est pas implémentée",
		},
		{
			Language:           "es-ES",
			DisplayKey:         "No Implementado",
			DisplayDescription: "Esta función no está implementada",
		},
		{
			Language:           "pt-BR",
			DisplayKey:         "Não Implementado",
			DisplayDescription: "Esta função não está implementada",
		},
		{
			Language:           "ru-RU",
			DisplayKey:         "Не Реализовано",
			DisplayDescription: "Эта функция не реализована",
		},
		{
			Language:           "de-DE",
			DisplayKey:         "Nicht Implementiert",
			DisplayDescription: "Diese Funktion ist nicht implementiert",
		},
		{
			Language:           "it-IT",
			DisplayKey:         "Non Implementato",
			DisplayDescription: "Questa funzione non è implementata",
		},
		{
			Language:           "nl-NL",
			DisplayKey:         "Niet Geïmplementeerd",
			DisplayDescription: "Deze functie is niet geïmplementeerd",
		},
		{
			Language:           "pl-PL",
			DisplayKey:         "Nie Zaimplementowano",
			DisplayDescription: "Ta funkcja nie jest zaimplementowana",
		},
		{
			Language:           "tr-TR",
			DisplayKey:         "Uygulanmadı",
			DisplayDescription: "Bu işlev uygulanmadı",
		},
		{
			Language:           "ar-SA",
			DisplayKey:         "غير مطبق",
			DisplayDescription: "هذه الوظيفة غير مطبقة",
		},
	},
	i18nErrHandlerNotFound: {
		{
			Language:           "en-US",
			DisplayDescription: "Handler [${handler}] for command path [${command}] is not found",
		},
		{
			Language:           "zh-CN",
			DisplayDescription: "未找到命令路径 [${command}] 的处理程序 [${handler}]",
		},
		{
			Language:           "ja-JP",
			DisplayDescription: "コマンドパス [${command}] のハンドラ [${handler}] が見つかりません",
		},
		{
			Language:           "ko-KR",
			DisplayDescription: "명령 경로 [${command}] 의 핸들러 [${handler}] 를 찾을 수 없습니다",
		},
		{
			Language:           "fr-FR",
			DisplayDescription: "Le gestionnaire [${handler}] pour le chemin de commande [${command}] n'est pas trouvé",
		},
		{
			Language:           "es-ES",
			DisplayDescription: "No se encuentra el controlador [${handler}] para la ruta de comando [${command}]",
		},
		{
			Language:           "pt-BR",
			DisplayDescription: "O manipulador [${handler}] para o caminho do comando [${command}] não foi encontrado",
		},
		{
			Language:           "ru-RU",
			DisplayDescription: "Обработчик [${handler}] для пути команды [${command}] не найден",
		},
		{
			Language:           "de-DE",
			DisplayDescription: "Handler [${handler}] für Befehlspfad [${command}] nicht gefunden",
		},
		{
			Language:           "it-IT",
			DisplayDescription: "Il gestore [${handler}] per il percorso del comando [${command}] non è stato trovato",
		},
		{
			Language:           "nl-NL",
			DisplayDescription: "Handler [${handler}] voor commandopad [${command}] niet gevonden",
		},
		{
			Language:           "pl-PL",
			DisplayDescription: "Nie znaleziono obsługi [${handler}] dla ścieżki komendy [${command}]",
		},
		{
			Language:           "tr-TR",
			DisplayDescription: "Komut yolu [${command}] için işleyici [${handler}] bulunamadı",
		},
		{
			Language:           "ar-SA",
			DisplayDescription: "لم يتم العثور على المعالج [${handler}] لمسار الأمر [${command}]",
		},
	},
	i18nErrInjectorNotFound: {
		{
			Language:           "en-US",
			DisplayDescription: "Prompt injector [${injector}] for command path [${command}] is not found",
		},
		{
			Language:           "zh-CN",
			DisplayDescription: "未找到命令路径 [${command}] 的提示注入器 [${injector}]",
		},
		{
			Language:           "ja-JP",
			DisplayDescription: "コマンドパス [${command}] のプロンプトインジェクタ [${injector}] が見つかりません",
		},
		{
			Language:           "ko-KR",
			DisplayDescription: "명령 경로 [${command}] 의 프롬프트 인젝터 [${injector}] 를 찾을 수 없습니다",
		},
		{
			Language:           "fr-FR",
			DisplayDescription: "Injecteur de prompt [${injector}] pour le chemin de commande [${command}] non trouvé",
		},
		{
			Language:           "es-ES",
			DisplayDescription: "El inyector de indicaciones [${injector}] para la ruta de comando [${command}] no se encuentra",
		},
		{
			Language:           "pt-BR",
			DisplayDescription: "O injetor de prompt [${injector}] para o caminho do comando [${command}] não foi encontrado",
		},
		{
			Language:           "ru-RU",
			DisplayDescription: "Инжектор приглашения [${injector}] для пути команды [${command}] не найден",
		},
		{
			Language:           "de-DE",
			DisplayDescription: "Prompt-Injektor [${injector}] für Befehlspfad [${command}] nicht gefunden",
		},
		{
			Language:           "it-IT",
			DisplayDescription: "L'iniettore di prompt [${injector}] per il percorso del comando [${command}] non è stato trovato",
		},
		{
			Language:           "nl-NL",
			DisplayDescription: "Prompt-injector [${injector}] voor commandopad [${command}] niet gevonden",
		},
		{
			Language:           "pl-PL",
			DisplayDescription: "Nie znaleziono wstrzykiwacza podpowiedzi [${injector}] dla ścieżki komendy [${command}]",
		},
		{
			Language:           "tr-TR",
			DisplayDescription: "Komut yolu [${command}] için ipucu enjektörü [${injector}] bulunamadı",
		},
		{
			Language:           "ar-SA",
			DisplayDescription: "لم يتم العثور على حقن الإشارة [${injector}] لمسار الأمر [${command}]",
		},
	},
	i18nErrPreferredTranslationNotConfigured: {
		{
			Language:           "en-US",
			DisplayDescription: "Description for command path [${command}] in preferred language [${language}] is not configured",
		},
		{
			Language:           "zh-CN",
			DisplayDescription: "未配置命令路径 [${command}] 在首选语言 [${language}] 的描述",
		},
		{
			Language:           "ja-JP",
			DisplayDescription: "優先言語 [${language}] のコマンドパス [${command}] の説明が構成されていません",
		},
		{
			Language:           "ko-KR",
			DisplayDescription: "선호하는 언어 [${language}] 에 대한 명령 경로 [${command}] 의 설명이 구성되지 않았습니다",
		},
		{
			Language:           "fr-FR",
			DisplayDescription: "Description du chemin de commande [${command}] dans la langue préférée [${language}] non configurée",
		},
		{
			Language:           "es-ES",
			DisplayDescription: "La descripción de la ruta de comando [${command}] en el idioma preferido [${language}] no está configurada",
		},
		{
			Language:           "pt-BR",
			DisplayDescription: "Descrição do caminho do comando [${command}] no idioma preferido [${language}] não configurada",
		},
		{
			Language:           "ru-RU",
			DisplayDescription: "Описание пути команды [${command}] на предпочитаемом языке [${language}] не настроено",
		},
		{
			Language:           "de-DE",
			DisplayDescription: "Beschreibung des Befehlspfads [${command}] in bevorzugter Sprache [${language}] nicht konfiguriert",
		},
		{
			Language:           "it-IT",
			DisplayDescription: "Descrizione del percorso del comando [${command}] nella lingua preferita [${language}] non configurata",
		},
		{
			Language:           "nl-NL",
			DisplayDescription: "Beschrijving van commandopad [${command}] in voorkeurstaal [${language}] niet geconfigureerd",
		},
		{
			Language:           "pl-PL",
			DisplayDescription: "Opis ścieżki komendy [${command}] w preferowanym języku [${language}] nie jest skonfigurowany",
		},
		{
			Language:           "tr-TR",
			DisplayDescription: "Tercih edilen dilde [${language}] komut yolu [${command}] için açıklama yapılandırılmamış",
		},
		{
			Language:           "ar-SA",
			DisplayDescription: "الوصف لمسار الأمر [${command}] في اللغة المفضلة [${language}] غير مكون",
		},
	},
}

var builtinLanguages = []string{
	"en-US", "zh-CN", "ja-JP", "ko-KR", "fr-FR", "es-ES", "pt-BR", "ru-RU", "de-DE", "it-IT", "nl-NL", "pl-PL", "tr-TR", "ar-SA",
}

func init() {
	i18nPacks = map[string]*TranslationSet{}
	for key, value := range fallbacks {
		ts := &TranslationSet{}
		ts.InitTranslations(value)
		i18nPacks[key] = ts
	}
}

func generateErrorPrompt(language string, key string, args map[string]string) []prompt.Suggest {
	name, description := i18nPacks[key].GetTranslation(language)
	description = values.NewStringTemplate(description, args).Parse()
	return []prompt.Suggest{{Text: name, Description: description}}
}

func BuiltinLanguages() []string {
	return builtinLanguages
}
