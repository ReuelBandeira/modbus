import 'dart:convert';

// import 'package:flutter_secure_storage/flutter_secure_storage.dart';
import 'package:http_interceptor/http/interceptor_contract.dart';
import 'package:http_interceptor/models/request_data.dart';
import 'package:http_interceptor/models/response_data.dart';
import 'package:isi_4_0/utils/navigator_service.dart';
import 'package:shared_preferences/shared_preferences.dart';

class ProtectedInterceptor implements InterceptorContract {
  // final secureStorage = const FlutterSecureStorage();

  Future<String> getToken() async {
    // String login = await secureStorage.read(key: 'login_details') ?? '';
    SharedPreferences preferences = await SharedPreferences.getInstance();
    String login = preferences.getString('login_details') ?? '';
    Map<String, dynamic> data = {};
    try {
      data = jsonDecode(login);
    } catch (e) {
      // readUser("isi40@grupoicts.com.br", "123456");
      // data = jsonDecode(login);
      data["token"] = '';
    }
    return data["token"];
  }

  @override
  Future<RequestData> interceptRequest({required RequestData data}) async {
    try {
      String token = await getToken();

      data.headers["Authorization"] = 'Bearer $token';
      data.headers["Content-Type"] = 'application/json; charset=UTF-8';
    } catch (e) {
      print(e);
    }
    return data;
  }

  @override
  Future<ResponseData> interceptResponse({required ResponseData data}) async {
    if (data.statusCode == 400 || data.statusCode == 401) {
      NavigationService().logOut();
    }

    return data;
  }
}
