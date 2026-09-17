import 'dart:convert';

import 'package:flutter/material.dart';
import 'package:isi_4_0/repository/db_api.dart';

class CardProvider extends ChangeNotifier {
  AuthService authService = AuthService.instance;

  String _configuredDevices = '0';
  String _usedMemory = '0%';
  String _diskUsed = '0%';
  String _connectedDevices = '0';

  String get configuredDevices => _configuredDevices;
  String get usedMemory => _usedMemory;
  String get diskUsed => _diskUsed;
  String get connectedDevices => _connectedDevices;

  Future<void> updateConfiguredDevices() async {
    _configuredDevices = await authService.countDevices();
    notifyListeners();
  }

  Future<void> updateUsedMemory() async {
    var jsonInfo = jsonDecode(await authService.getInfoUnits());
    if (jsonInfo.runtimeType.toString() == "List<dynamic>") {
      _usedMemory = "${jsonInfo[0]['memoryUsage'] ?? "0"}%";
      notifyListeners();
    }
  }

  Future<void> updateDiskUsed() async {
    var jsonInfo = jsonDecode(await authService.getInfoUnits());
    if (jsonInfo.runtimeType.toString() == "List<dynamic>") {
      _diskUsed = "${jsonInfo[0]['diskUsage'] ?? "0"}%";
      notifyListeners();
    }
  }

  Future<void> updateConnectedDevices() async {
    var jsonInfo = jsonDecode(await authService.getInfoUnits());
    if (jsonInfo.runtimeType.toString() == "List<dynamic>") {
      _connectedDevices = jsonInfo[0]['connectedDevices'].length.toString();
      notifyListeners();
    }
  }

  void updateAll() {
    updateConfiguredDevices();
    updateUsedMemory();
    updateDiskUsed();
    updateConnectedDevices();
  }
}
