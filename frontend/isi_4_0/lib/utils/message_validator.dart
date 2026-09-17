import 'package:isi_4_0/utils/enum_login.dart';
import 'package:isi_4_0/utils/mask.dart';

class MessageValidator {
  MessageValidator._();

  static final MessageValidator _instance = MessageValidator._();

  static MessageValidator get instance => _instance;

  String? getMessage(CurrentPage currentPage, String? value,
      CurrentTextField currentTextField, messageError) {
    switch (currentPage) {
      case CurrentPage.signIn:
        switch (currentTextField) {
          case CurrentTextField.email:
            if (value!.isEmpty) {
              return '';//messageError['msg_email_empty'];
            } else {
              if (MaskDetector.detectMaskType(value) == "Unknown") {
                return '';//messageError['msg_email_wrong'];
              }
            }
            break;
          case CurrentTextField.password:
            if (value!.isEmpty) {
              return '';//messageError['msg_password_empty'];
            }
            break;
          case CurrentTextField.user:
            break;
          case CurrentTextField.confirmPassword:
            break;
        }
        break;
      case CurrentPage.recover:
        switch (currentTextField) {
          case CurrentTextField.email:
            if (value!.isEmpty) {
              return '';//messageError['msg_email_empty'];
            } else {
              if (MaskDetector.detectMaskType(value) == "Unknown") {
                return '';//messageError['msg_email_wrong'];
              }
            }
            break;
          case CurrentTextField.password:
            if (value!.isEmpty) {
              return '';//messageError['msg_password_empty'];
            }
            break;
          case CurrentTextField.confirmPassword:
            if (value!.isEmpty) {
              return '';//messageError['msg_confirm_password_empty'];
            }
            break;
          case CurrentTextField.user:
            break;
        }
        break;
      case CurrentPage.register:
        switch (currentTextField) {
          case CurrentTextField.email:
            if (value!.isEmpty) {
              return '';//messageError['msg_email_empty'];
            } else {
              if (MaskDetector.detectMaskType(value) == "Unknown") {
                return '';//messageError['msg_email_wrong'];
              }
            }
            break;
          case CurrentTextField.password:
            if (value!.isEmpty) {
              return '';//messageError['msg_password_empty'];
            }
            break;
          case CurrentTextField.confirmPassword:
            if (value!.isEmpty) {
              return '';//messageError['msg_confirm_password_empty'];
            }
            break;
          case CurrentTextField.user:
            if (value!.isEmpty || value.trim() == "") {
              return '';//messageError['msg_user_empty'];
            } else {
              if (value.length < 3) {
                return '';//messageError['msg_low_user'];
              }
            }

            break;
        }
        break;
      case CurrentPage.settings:
        switch (currentTextField) {
          case CurrentTextField.addressIP:
            if (value!.isEmpty) {
              return '';//messageError['devices']['msg_alert_field'];
            } else {
              if ((MaskDetector.detectMaskType(value) == "IP") != true) {
                return '';//messageError['devices']['msg_ip_wrong'];
              }
            }
            break;
          case CurrentTextField.port:
            if (value!.isEmpty) {
              return '';//messageError['devices']['msg_alert_field'];
            }
            break;
          case CurrentTextField.user:
            break;
          case CurrentTextField.password:
            break;
        }
        break;
      default:
    }
    return null;
  }
}
