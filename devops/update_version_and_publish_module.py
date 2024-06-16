
import os
import yaml

configYmlPath = './service.yml'
print('configYmlPath: ', configYmlPath)
configYml = open(configYmlPath)
config = yaml.safe_load(configYml)
lastVersionKey = 'last_version'
lastVersion = config[lastVersionKey]
print('last version: ' + lastVersion)
v1, v2, v3 = str(lastVersion).split('.')
v3 = str(int(v3) + 1)
noVprefixBuildVersion = v1 + '.' + v2 + '.' + v3
buildVersion = 'v' + noVprefixBuildVersion
print('build version: ' + buildVersion)
config[lastVersionKey] = noVprefixBuildVersion
with open(configYmlPath, 'w') as updateYmlPath:
	yaml.dump(config, updateYmlPath)

os.system('./devops/publish_module.sh ' + buildVersion)