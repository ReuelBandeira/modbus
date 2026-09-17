import 'dart:convert' as convert;
import 'dart:convert';
import 'dart:html' as html;

// import 'package:flutter_secure_storage/flutter_secure_storage.dart';
import 'package:http/http.dart' as http;
import 'package:http_interceptor/http/intercepted_http.dart';
import 'package:isi_4_0/repository/protected_interceptor.dart';
import 'package:isi_4_0/utils/enum_db.dart';
import 'package:shared_preferences/shared_preferences.dart';

APIConfig apiConfig = APIConfig();

class AuthService {
  static final AuthService instance = AuthService._internal();
  static final InterceptedHttp protectedHttp = InterceptedHttp.build(interceptors: [
    ProtectedInterceptor(),
  ]);
  // final secureStorage = const FlutterSecureStorage();

  factory AuthService() {
    return instance;
  }

  AuthService._internal();

  Future<void> refreshPage(String page) async {
    SharedPreferences prefs = await SharedPreferences.getInstance();
    await prefs.setString('initial_page', page);
  }

  Future<String> getInitialPage() async {
    SharedPreferences prefs = await SharedPreferences.getInstance();
    return prefs.getString('initial_page') ?? '/';
  }

  Future<void> saveLang(String lang) async {
    SharedPreferences prefs = await SharedPreferences.getInstance();
    await prefs.setString('choose_lang', lang);
  }

  Future<String> getLang() async {
    SharedPreferences prefs = await SharedPreferences.getInstance();
    String? userLang = prefs.getString('choose_lang');
    userLang ??= 'pt';//await getLanguage();
    return userLang;
  }

  Future<void> saveLoginDetails(String data) async {
    // await secureStorage.write(key: 'login_details', value: data);
    SharedPreferences login = await SharedPreferences.getInstance();
    await login.setString('login_details', data);
  }

  Future<String> getLoginDetails() async {
    // String data = await secureStorage.read(key: 'login_details') ?? '';
    SharedPreferences login = await SharedPreferences.getInstance();
    String data = login.getString('login_details') ?? '';
    return data;
  }

  Future<void> removeLogin() async {
    SharedPreferences login = await SharedPreferences.getInstance();
    await login.remove('login_details');
    // await secureStorage.delete(key: 'login_details');
  }

  Future<String> getToken() async {
    String login = await getLoginDetails();
    late Map<String, dynamic> data;
    try {
      data = convert.jsonDecode(login);
    } catch (e) {
      data["token"] = '';
    }
    return data["token"];
  }

  Future<int> readUser(String email, String password) async {
    int status = 503;
    try {
      final url = Uri.parse(apiConfig.url(APIMethod.loginRead.id));
      final response = await http.post(
        url,
        headers: <String, String>{
          'Content-Type': 'application/json; charset=UTF-8',
        },
        body: convert.jsonEncode({"email": email, "password": password}),
      );
      if (response.statusCode == 200) {
        final decodedBody = utf8.decode(response.bodyBytes);
        await saveLoginDetails(decodedBody);
      }
      status = response.statusCode;
    } catch (e) {
      e.toString();
    }
    return status;
  }

  // Creates
  Future<int> createUser(
      String name, String email, String password, String language) async {
    int status = 503;
    try {
      final url = Uri.parse(apiConfig.url(APIMethod.loginCreate.id));
      Map<String, String> headers = {
        'Content-Type': 'application/json'
      };
      final response = await http.post(
        url,
        headers: headers,
        body: convert.jsonEncode({
          "name": name,
          "email": email,
          "password": password,
          "language": language
        }),
      );
      status = response.statusCode;
    } catch (e) {
      e.toString();
    }
    return status;
  }

  Future<String> countDevices() async {
    String count = "0";
    try {
      final url = Uri.parse(apiConfig.url(APIMethod.countDevice.id));
      final response = await protectedHttp.get(
        url,
      );

      if (response.statusCode == 200) {
        Map<String, dynamic> jsonData = json.decode(response.body);
        count = jsonData['count'].toString();
      }
    } catch (e) {
      e.toString();
    }
    return count;
  }

  Future<Map<String, dynamic>> getDevices(int page, int items) async {
    Map<String, dynamic> data = {};
    try {
      final url = Uri.parse(apiConfig.url(APIMethod.getDevices.id));
      final response = await protectedHttp.get(
        url,
        params: {
          'page': page,
          'items': items
        }
      );
      if (response.statusCode == 200) {
        Map<String, dynamic> jsonData = jsonDecode(response.body);
        data = jsonData;
      }
    } catch (e) {
      e.toString();
    }
    return data;
  }

  Future<String> getInfoUnits() async {
    dynamic data;
    try {
      final url = Uri.parse(apiConfig.url(APIMethod.infoUnits.id));
      final response = await protectedHttp.get(
        url,
      );
      if (response.statusCode == 200) {
        data = response.body;
        if (data == "null") {
          data = '[{"memoryUsage":0,"diskUsage":0,"connectedDevices":""}]';
        }
      }
    } catch (e) {
      e.toString();
    }
    return data;
  }

  Future<int> saveProtocolMQTT(
      String server,
      String port,
      String username,
      String password,
      String inputInfo,
      String errorInfo,
      String statusDevice) async {
    int status = 503;
    try {
      final url = Uri.parse(apiConfig.url(APIMethod.saveMQTT.id));
      final response = await protectedHttp.post(
        url,
        body: convert.jsonEncode({
          "server": server,
          "port": port,
          "username": username,
          "password": password,
          "inputInfo": inputInfo,
          "errorInfo": errorInfo,
          "statusDevice": statusDevice
        }),
      );
      status = response.statusCode;
    } catch (e) {
      e.toString();
    }
    return status;
  }

  Future<http.Response> saveDevice(String device) async {
    http.Response responseReturn = http.Response('', 503);
    try {
      final url = Uri.parse(apiConfig.url(APIMethod.saveDevice.id));
      final response = await protectedHttp.post(
        url,
        body: device,
      );
      responseReturn = response;
    } catch (e) {
      e.toString();
    }
    return responseReturn;
  }

  Future<bool> deleteDevice(String deviceID) async {
    var status = false;
    try {
      final url =
          Uri.parse(apiConfig.url(APIMethod.deleteDevices.id) + '/$deviceID');
      final response = await protectedHttp.delete(
        url,
      );
      response.statusCode == 200 ? status = true : status = false;
    } catch (e) {
      e.toString();
    }
    return status;
  }

  Future<int> updateDevice(String device, String data) async {
    int status = 503;
    try {
      final url =
          Uri.parse(apiConfig.url(APIMethod.updateDevices.id) + '/$device');
      final response = await protectedHttp.patch(
        url,
        body: data,
      );
      status = response.statusCode;
    } catch (e) {
      e.toString();
    }
    return status;
  }

  Future<int> updateUser(String email, String password) async {
    int status = 503;
    try {
      final url = Uri.parse(apiConfig.url(APIMethod.loginUpdate.id) + email);
      final response = await protectedHttp.put(
        url,
        body: convert.jsonEncode(
            {"name": "needToRemove", "email": email, "password": password}),
      );
      status = response.statusCode;
    } catch (e) {
      e.toString();
    }
    return status;
  }

  Future<int> updateProfile(String id, String name, String email,
      String password, String language) async {
    int status = 503;
    String token = await getToken();
    try {
      final url = Uri.parse(apiConfig.url(APIMethod.updateProfile.id) + id);
      final response = await protectedHttp.patch(
        url,
        body: convert.jsonEncode({
          "email": email,
          "name": name,
          "password": password,
          "language": language
        }),
      );
      status = response.statusCode;
      if (status == 200) {
        String data = '{"email":"$email",'
                        '"id":$id,'
                        '"language":"$language",'
                        '"name":"$name",'
                        '"token":"$token"}';
        await saveLoginDetails(data);
      }
    } catch (e) {
      e.toString();
    }
    return status;
  }

  Future<List<dynamic>> getProtocols() async {
    dynamic data;
    try {
      final url = Uri.parse(apiConfig.url(APIMethod.getProtocols.id));

      final response = await protectedHttp.get(
        url,
      );

      if (response.statusCode == 200) {
        List<dynamic> jsonData = json.decode(response.body);
        data = jsonData;
      } else {
        data = "";
      }
    } catch (e) {
      e.toString();
    }
    return data;
  }

  Future<String> getLanguage() async {
    String data = "pt";
    String token = await getToken();

    try {
      final url = Uri.parse(apiConfig.url(APIMethod.getLanguage.id));

      final response = await http.get(
        url,
        headers: <String, String>{
          'Authorization': 'Bearer $token',
          'Content-Type': 'application/json; charset=UTF-8',
        },
      );

      if (response.statusCode == 200) {
        List<dynamic> jsonData = json.decode(response.body);
        data = jsonData[0]["language"];
      }
    } catch (e) {
      print(e);
    }
    return data;
  }

  Future<Map<String, dynamic>> getMqttStatus() async {
    Map<String, dynamic> data = {};
    try {
      final url = Uri.parse(apiConfig.url(APIMethod.getMqttStatus.id));

      final response = await protectedHttp.get(
        url,
      );

      if (response.statusCode == 200) {
        Map<String, dynamic> jsonMap = json.decode(response.body);
        data = jsonMap;
      }
    } catch (e) {
      e.toString();
    }
    return data;
  }

  Future<Map<String, dynamic>> getLogList(int page, int items) async {
    Map<String, dynamic> data = {};
    try {
      final url = Uri.parse(apiConfig.url(APIMethod.getLogList.id));

      final response = await protectedHttp.get(
        url,
        params: {
          'page': page,
          'items': items
        }
      );

      if (response.statusCode == 200) {
        Map<String, dynamic> jsonData = jsonDecode(response.body);
        data = jsonData;
      }
    } catch (e) {
      e.toString();
    }
    return data;
  }

  Future<void> downloadLogFile(file) async {
    try {
      final anchor = html.AnchorElement(
          href: apiConfig.url(APIMethod.downloadLogFile.id) + file.toString())
        ..target = 'blank'
        ..click();
      anchor.remove();
    } catch (e) {
      rethrow;
    }
  }

  Future<void> downloadAllLogFiles() async {
    try {
      final anchor = html.AnchorElement(
          href: apiConfig.url(APIMethod.downloadAllLogFiles.id))
        ..target = 'blank'
        ..click();
      anchor.remove();
    } catch (e) {
      rethrow;
    }
  }
}
