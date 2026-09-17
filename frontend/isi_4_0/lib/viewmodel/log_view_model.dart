import 'dart:js_util';

import 'package:flutter/material.dart';
import 'package:isi_4_0/repository/db_api.dart';

class LogViewModel extends ChangeNotifier {
  AuthService authService = AuthService.instance;

  Future<Map<String, dynamic>> getLogList(int page, int items) {
    return authService.getLogList(page, items);
  }

  Future<void> downloadLogFile(file) async {
    try {
      return await authService.downloadLogFile(file);
    } catch (e) {
      rethrow;
    }
  }

  Future<void> downloadAllLogFile() async {
    try {
      return await authService.downloadAllLogFiles();
    } catch (e) {
      rethrow;
    }
  }
}
