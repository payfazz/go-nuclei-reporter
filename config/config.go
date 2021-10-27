package config

var base = mergeConfig(
	fileLocationConfig,
	logLevelConfig,
	githubTokenConfig,
	organizationNameConfig,
	notionTokenConfig,
	notionReportDatabaseConfig,
	notionChangeDatabaseConfig,
)
