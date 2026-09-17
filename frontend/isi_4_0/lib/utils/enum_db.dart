class APIConfig {
  String hostServer = getUrl(Uri.base.origin);
  String portServer = '8585';

  APIConfig();

  url(int method) {
    switch (APIMethod.values[method]) {
      case APIMethod.loginRead:
        return hostServer + portServer + DBRoutes.loginReadUser.id;
      case APIMethod.loginCreate:
        return hostServer + portServer + DBRoutes.loginCreateUser.id;
      case APIMethod.loginUpdate:
        return hostServer + portServer + DBRoutes.loginUpdateUser.id;
      case APIMethod.saveMQTT:
        return hostServer + portServer + DBRoutes.saveMQTTConfig.id;
      case APIMethod.saveDevice:
        return hostServer + portServer + DBRoutes.saveDeviceConfig.id;
      case APIMethod.countDevice:
        return hostServer + portServer + DBRoutes.countDeviceConfig.id;
      case APIMethod.getDevices:
        return hostServer + portServer + DBRoutes.getDeviceConfig.id;
      case APIMethod.deleteDevices:
        return hostServer + portServer + DBRoutes.deleteDeviceConfig.id;
      case APIMethod.updateDevices:
        return hostServer + portServer + DBRoutes.updateDeviceConfig.id;
      case APIMethod.infoUnits:
        return hostServer + portServer + DBRoutes.readInfoUnits.id;
      case APIMethod.getProtocols:
        return hostServer + portServer + DBRoutes.readProtocols.id;
      case APIMethod.updateProfile:
        return hostServer + portServer + DBRoutes.updateProfileConfig.id;
      case APIMethod.getLanguage:
        return hostServer + portServer + DBRoutes.getLanguageDefault.id;
      case APIMethod.getMqttStatus:
        return hostServer + portServer + DBRoutes.getMqttStatus.id;
      case APIMethod.getLogList:
        return hostServer + portServer + DBRoutes.getLogList.id;
      case APIMethod.downloadLogFile:
        return hostServer + portServer + DBRoutes.downloadLogFile.id;
      case APIMethod.downloadAllLogFiles:
        return hostServer + portServer + DBRoutes.downloadAllLogFile.id;
      default:
    }
  }
}

enum APIMethod<int> {
  loginRead(id: 0),
  loginCreate(id: 1),
  loginUpdate(id: 2),
  saveMQTT(id: 3),
  saveDevice(id: 4),
  countDevice(id: 5),
  getDevices(id: 6),
  deleteDevices(id: 7),
  updateDevices(id: 8),
  infoUnits(id: 9),
  getProtocols(id: 10),
  updateProfile(id: 11),
  getLanguage(id: 12),
  getMqttStatus(id: 13),
  getLogList(id: 14),
  downloadLogFile(id: 15),
  downloadAllLogFiles(id: 16);

  final int id;

  const APIMethod({required this.id});
}

enum DBRoutes {
  loginReadUser(id: '/api/v1/login/'),
  loginCreateUser(id: '/api/v1/register/'),
  loginUpdateUser(id: '/api/v1/user/'),
  saveMQTTConfig(id: '/api/v1/mqtt'),
  saveDeviceConfig(id: '/api/v1/configurations'),
  countDeviceConfig(id: '/api/v1/configurationcount'),
  getDeviceConfig(id: '/api/v1/configurations'),
  deleteDeviceConfig(id: '/api/v1/configurations'),
  updateDeviceConfig(id: '/api/v1/configurations'),
  readInfoUnits(id: '/api/v1/deviceinfo'),
  readProtocols(id: '/api/v1/protocols'),
  updateProfileConfig(id: '/api/v1/userdata/'),
  getLanguageDefault(id: '/api/v1/language'),
  getMqttStatus(id: '/api/v1/mqttstatus'),
  getLogList(id: '/api/v1/logs'),
  downloadLogFile(id: '/api/v1/logfile/'),
  downloadAllLogFile(id: '/api/v1/logfiles');

  final String id;

  const DBRoutes({required this.id});
}

String getUrl(String url) {
  int indexPrimeiroDoisPontos = url.indexOf(':');
  if (indexPrimeiroDoisPontos != -1) {
    int indexSegundoDoisPontos = url.indexOf(':', indexPrimeiroDoisPontos + 1);
    if (indexSegundoDoisPontos != -1) {
      return url.substring(0, indexSegundoDoisPontos + 1);
    }
  }
  return "http://localhost:";
}
