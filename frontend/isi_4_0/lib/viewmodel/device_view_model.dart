import 'dart:convert';

import 'package:flutter/material.dart';
import 'package:isi_4_0/model/cards_model.dart';
import 'package:isi_4_0/pages/devices_extra_data.dart';
import 'package:isi_4_0/providers/lang.dart';
import 'package:isi_4_0/repository/db_api.dart';
import 'package:isi_4_0/utils/data_devices.dart';
import 'package:loggy/loggy.dart';

class DeviceViewModel extends ChangeNotifier {
  AuthService authService = AuthService();

  DataDevices _dataDevices = DataDevices(id: '', fields: [], topics: []);
  DataDevices get dataDevices => _dataDevices;

  Future<void> updateDataDevices(DataDevices devices) async {
    _dataDevices = devices;
    notifyListeners();
  }

  Future<void> devicesExtraData(BuildContext context) async {
    Navigator.push(
      context,
      MaterialPageRoute(
        builder: (BuildContext context) =>
            DevicesExtraData(dataDevices: _dataDevices),
      ),
    );
  }

  Future<List<InfoCards>> getDataCards(LanguageProvider language) async {
    List<InfoCards> infoCards = [];
    try {
      String dataCard = await authService.getInfoUnits();
      var lang =
          language.getDataLanguage(language.currentLanguage)['cardTitle'];
      Map<String, dynamic> data = json.decode(dataCard);
      data.forEach((key, value) {
        if (key == 'count') {
          infoCards.add(InfoCards(lang[key], value.toString()));
        } else {
          infoCards.add(InfoCards(lang[key], '${value.toStringAsFixed(2)}%'));
        }
      });
    } catch (error) {
      logError('Error fetching data on function "getDataCards: $error');
    }
    return infoCards;
  }

  Future<Map<String, dynamic>> getDataDevices(int page, int items) async {
    var getDevices = await authService.getDevices(page, items);
    return getDevices;
  }
}
