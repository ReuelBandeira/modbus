import 'package:flutter/services.dart';

class CustomErrorFilterInput extends TextInputFormatter {
  final Function(bool, String) onError;
  final bool? isRequired;
  final bool? isEmail;
  final bool? isUser;
  final int? min;
  final int? max;
  final bool? digitsOnly;
  final bool? ipOnly;
  final bool? enabled;
  final Map<String, dynamic> translator;
  CustomErrorFilterInput(
      {
        this.isRequired,
        this.isEmail,
        this.isUser,
        this.min,
        this.max,
        this.digitsOnly,
        this.ipOnly,
        this.enabled,
        required this.onError,
        required this.translator,
      }
  );

  @override
  TextEditingValue formatEditUpdate(TextEditingValue oldValue, TextEditingValue newValue) {
    String errorString = '';
    bool match = false;

    if (enabled != null && !enabled!) {
      onError(match, errorString);
      return newValue;
    }

    if (isRequired != null && isRequired == true) {
      bool matchReq = (newValue.text == '');
      if (matchReq) {
        errorString = translator["customTextField"]["errorRequired"];
        onError(matchReq, errorString);
        return newValue;
      }
    }

    if (isEmail != null && isEmail == true) {
      match = !RegExp(r'^[\w-\.]+@([\w-]+\.)+[\w-]{2,4}$').hasMatch(newValue.text);
      if (match) {
        errorString = translator["dictionary"]["msg_email_wrong"];
      }
    }

    if (isUser != null && isUser == true) {
      if (newValue.text.trim() == "") {
        match = true;
        errorString = translator["dictionary"]["msg_user_empty"];
      } else {
        if (newValue.text.length < 3) {
          match = true;
          errorString = translator["dictionary"]["msg_low_user"];
        }
      }
    }

    if (digitsOnly != null && digitsOnly == true) {
      match = RegExp(r'[^0-9]').hasMatch(newValue.text);
      if (match) {
        errorString = translator["customTextField"]["errorDigitsOnly"];
      }
    }

    if (min != null) {
      if (digitsOnly != null && digitsOnly == true) {
        int? val = int.tryParse(newValue.text);
        if (val == null) {
          double? valDouble = double.tryParse(newValue.text);
          if (valDouble != null) {
            if (valDouble < min!) {
              errorString = translator["customTextField"]["errorMinReq"] + min.toString();
              match = true;
            }
          }
        } else {
          if (val < min!) {
            errorString = translator["customTextField"]["errorMinReq"] + min.toString();
            match = true;
          }
        }
      } else {
        if (newValue.text.length < min!) {
          errorString = translator["customTextField"]["errorMinStringReq"] +
              min.toString() + translator["customTextField"]["labelCharacters"];
          match = true;
        }
      }
    }

    if (max != null) {
      if (digitsOnly != null && digitsOnly == true) {
        int? val = int.tryParse(newValue.text);
        if (val == null) {
          double? valDouble = double.tryParse(newValue.text);
          if (valDouble != null) {
            if (valDouble > max!) {
              errorString = translator["customTextField"]["errorMaxReq"] + max.toString();
              match = true;
            }
          }
        } else {
          if (val > max!) {
            errorString = translator["customTextField"]["errorMaxReq"] + max.toString();
            match = true;
          }
        }
      } else {
        if (newValue.text.length > max!) {
          errorString = translator["customTextField"]["errorMaxStringReq"] +
              max.toString() + translator["customTextField"]["labelCharacters"];
          match = true;
        }
      }
    }

    if (ipOnly != null && ipOnly == true) {
      // Regex for valid ip address only
      // match = !RegExp(r'^(?:(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}'
      // r'(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)$').hasMatch(newValue.text);
      match = !RegExp(r'^[\.0-9]*$').hasMatch(newValue.text);
      if (match) {
        errorString = translator["customTextField"]["errorIp"];
      }
    }

    onError(match, errorString);
    return newValue;
  }

  void checkErrors(String initialText) {
    TextEditingValue newText = TextEditingValue(text: initialText);
    formatEditUpdate(newText, newText);
  }
}
