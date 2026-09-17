import 'dart:math';

import 'package:flutter/material.dart';

class CustomColors {
  static const Color background800 = Color(0xFFFBFBFB);
  static const Color background700 = Color(0xFFFFFFFF);
  static const Color background600 = Color(0xFF121212);

  static const Color success800 = Color(0xFF28731E);
  static const Color success700 = Color(0xFF379031);
  static const Color success600 = Color(0xFF4AA849);
  static const Color success500 = Color(0xFF66BB6A);
  static const Color success400 = Color(0xFF8DCF8D);
  static const Color success300 = Color(0xFFB8E2B5);
  static const Color success200 = Color(0xFFE1F3DF);

  static const Color info800 = Color(0xFF008080);
  static const Color info700 = Color(0xFF00A4AA);
  static const Color info600 = Color(0xFF06AFD5);
  static const Color info500 = Color(0xFF29B6F6);
  static const Color info400 = Color(0xFF5DD4FD);
  static const Color info300 = Color(0xFF95EBFF);
  static const Color info200 = Color(0xFFD0FAFF);

  static const Color error800 = Color(0xFF4D0000);
  static const Color error700 = Color(0xFF880002);
  static const Color error600 = Color(0xFFC41411);
  static const Color error500 = Color(0xFFF44336);
  static const Color error400 = Color(0xFFFF4F4D);
  static const Color error300 = Color(0xFFFF7779);
  static const Color error200 = Color(0xFFFFB3B3);

  static const Color error600Fade = Color(0x77C41411);

  static const Color alert800 = Color(0xFF806000);
  static const Color alert700 = Color(0xFFAA8900);
  static const Color alert600 = Color(0xFFD5AF0F);
  static const Color alert500 = Color(0xFFFFCF34);
  static const Color alert400 = Color(0xFFFFD665);
  static const Color alert300 = Color(0xFFFFE09A);
  static const Color alert200 = Color(0xFFFFEFD2);

  static const Color neutral800 = Color(0xFF2E2E2E);
  static const Color neutral700 = Color(0xFF525252);
  static const Color neutral600 = Color(0xFF767676);
  static const Color neutral500 = Color(0xFF9A9A9A);
  static const Color neutral400 = Color(0xFFB8B8B8);
  static const Color neutral300 = Color(0xFFD5D5D5);
  static const Color neutral200 = Color(0xFFF3F3F3);

  static const Color secondary200 = Color(0xFFFFBFB3);
  static const Color secondary300 = Color(0xFFFF9E77);
  static const Color secondary400 = Color(0xFFF8853C);
  static const Color secondary500 = Color(0xFFD67018);
  static const Color secondary600 = Color(0xFFBC4900);
  static const Color secondary700 = Color(0xFF882700);
  static const Color secondary800 = Color(0xFF4D0D00);

  static const Color primary200 = Color(0xFFB3D5E1);
  static const Color primary300 = Color(0xFF7AADC2);
  static const Color primary400 = Color(0xFF4A809D);
  static const Color primary500 = Color(0xFF205171);
  static const Color primary600 = Color(0xFF0E4561);
  static const Color primary700 = Color(0xFF03364B);
  static const Color primary800 = Color(0xFF00232E);

  static const Color neutral700Faded = Color(0x66525252);
  static const Color dropShadow = Color(0x29000000);

  static const Color primaryColorApp = Color.fromRGBO(32, 81, 113, 1);
  static const Color whiteColorLow = Color.fromARGB(255, 249, 251, 255);
  static const Color whiteColorHigh = Color.fromRGBO(243, 243, 243, 1);
  static const Color greyColorHigh = Color.fromRGBO(82, 82, 82, 1);
  static const Color blueColorLow = Color.fromRGBO(130, 210, 237, 1);

  static const Color secondaryColorApp = Color.fromRGBO(29, 90, 113, 1);

  static const Color greenColorLight = Color.fromRGBO(161, 219, 233, 0.282);

  static const Color blueColorHigh = Color.fromRGBO(17, 38, 60, 1);
  static const Color yellowColorLight = Color.fromRGBO(253, 207, 111, 1);

  static const Color defaultColorApp = Color.fromARGB(0, 255, 255, 255);
  static const Color blackColorAppLow = Colors.black54;
  static const Color backgroundTextFieldDropButton = Color.fromRGBO(
      193, 193, 193, 0.3);
  static const Color greyDefault = Colors.grey;
  static const Color backgroundTextFieldProtocols = Color.fromARGB(255, 96, 144, 177);
  static const Color backgroundPages = Color.fromRGBO(209, 235, 244, 1);
  static const Color backgroundSquareHome = Color.fromARGB(33, 18, 35, 37);

  static const Color hoverColorLogin = Color.fromARGB(58, 176, 174, 174);

  static const Color lineColors = Color.fromARGB(12, 0, 0, 0);

  static const Color deviceOnline = info500;
  static const Color deviceOffOnline = error600;
}

Color getRandomColor() {
  final random = Random();
  final r = random.nextInt(256);
  final g = random.nextInt(256);
  final b = random.nextInt(256);
  final o = random.nextDouble();
  return Color.fromRGBO(r, g, b, o);
}
