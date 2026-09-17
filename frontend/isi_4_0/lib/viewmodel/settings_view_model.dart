import 'package:flutter/material.dart';
import 'package:isi_4_0/repository/db_api.dart';

class SettignsViewModel extends ChangeNotifier {
  AuthService authService = AuthService.instance;

  Future<int> saveMQTT() async {
    return 0;
  }

  Future<Map<String, dynamic>> getMqttStatus() async {
    return await authService.getMqttStatus();
  }
}
