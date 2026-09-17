import 'package:flutter/material.dart';

class ProfileProvider extends ChangeNotifier {
  bool _name = false;
  bool get currentName => _name;

  bool _email = false;
  bool get currentEmail => _email;

  bool _password = false;
  bool get currentPassword => _password;

  void updateNameState(bool state) {
    _name = state;
    notifyListeners();
  }

  void updateEmailState(bool state) {
    _email = state;
    notifyListeners();
  }

  void updatePasswordState(bool state) {
    _password = state;
    notifyListeners();
  }
}
