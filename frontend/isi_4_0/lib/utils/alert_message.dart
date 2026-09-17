import 'package:flutter/material.dart';
import 'package:fluttertoast/fluttertoast.dart';
import 'package:toastification/toastification.dart';

class AlertMessage {
  AlertMessage();
  message(String msg, Color appBgColor, String webBgColor) {
    Fluttertoast.showToast(
        msg: msg,
        toastLength: Toast.LENGTH_SHORT,
        gravity: ToastGravity.BOTTOM,
        timeInSecForIosWeb: 3,
        backgroundColor: appBgColor,
        textColor: Colors.white,
        webPosition: 'right',
        webBgColor: "linear-gradient(to right, $webBgColor, $webBgColor)",
        fontSize: 16.0);
  }

  serverMessage(int statusCode, String msg) {
    late Color appBgColor;
    late String webBgColor;

    if (statusCode > 0 && statusCode < 200) {
      appBgColor = Colors.blue;
      webBgColor = '#4B7EFF';
    }
    if (statusCode >= 200 && statusCode < 300) {
      appBgColor = Colors.green;
      webBgColor = '#4caf50';
    }
    if (statusCode >= 300 && statusCode < 400) {
      appBgColor = Colors.orange;
      webBgColor = '#FF9800';
    }
    if (statusCode >= 400 && statusCode < 500) {
      appBgColor = Colors.red;
      webBgColor = '#dc1c13';
    }
    if (statusCode >= 500) {
      appBgColor = Colors.black;
      webBgColor = '#1C1D1C';
    }

    Fluttertoast.showToast(
        msg: msg,
        toastLength: Toast.LENGTH_SHORT,
        gravity: ToastGravity.TOP,
        timeInSecForIosWeb: 3,
        backgroundColor: appBgColor,
        textColor: Colors.white,
        webPosition: 'right',
        webBgColor: "linear-gradient(to right, $webBgColor, $webBgColor)",
        fontSize: 16.0);
  }

  showError(BuildContext context, String msg) {
    toastification.show(
      context: context,
      type: ToastificationType.error,
      style: ToastificationStyle.fillColored,
      showProgressBar: false,
      closeButtonShowType: CloseButtonShowType.onHover,
      closeOnClick: true,
      title: msg,
      autoCloseDuration: const Duration(seconds: 3),
      // margin: const EdgeInsets.only(top: 35, bottom: 0, left: 12, right: 12),
    );
  }

  showSuccess(BuildContext context, String msg) {
    toastification.show(
      context: context,
      type: ToastificationType.success,
      style: ToastificationStyle.fillColored,
      showProgressBar: false,
      closeButtonShowType: CloseButtonShowType.onHover,
      closeOnClick: true,
      title: msg,
      autoCloseDuration: const Duration(seconds: 3),
      // margin: const EdgeInsets.only(top: 35, bottom: 0, left: 12, right: 12),
    );
  }
}
