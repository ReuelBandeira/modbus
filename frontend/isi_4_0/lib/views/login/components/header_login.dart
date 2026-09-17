import 'package:flutter/material.dart';
import 'package:isi_4_0/providers/lang.dart';
import 'package:isi_4_0/utils/app_details.dart';

class HeaderLogin extends StatefulWidget {
  final double height;
  final double width;
  final LanguageProvider languageProvider;
  const HeaderLogin(
      {super.key,
      required this.height,
      required this.width,
      required this.languageProvider});
  @override
  State<HeaderLogin> createState() => _HeaderLoginState();
}

class _HeaderLoginState extends State<HeaderLogin> {
  final double sizeFlag = 24;
  @override
  Widget build(BuildContext context) {
    return Container(
      height: widget.height,
      width: widget.width,
      alignment: Alignment.centerRight,
      margin: const EdgeInsets.only(right: 4),
      child: Container(
        height: widget.height,
        width: widget.height,
        alignment: Alignment.center,
        padding: const EdgeInsets.symmetric(horizontal: 8.0),
        child: InkWell(
          child: Image.asset(
            getCurrentFlag(widget.languageProvider.currentLanguage),
            height: sizeFlag,
            width: sizeFlag,
          ),
          onTap: () {
            showMenu(
              context: context,
              position: const RelativeRect.fromLTRB(60, 44, 16, 0),
              items: List.generate(
                flags.length,
                (flag) {
                  return PopupMenuItem(
                    child: Image.asset(
                      flags[flag]['image'],
                      height: sizeFlag,
                      width: sizeFlag,
                    ),
                    onTap: () {
                      widget.languageProvider
                          .updateCurrentLanguage(flags[flag]['lang']);
                    },
                  );
                },
              ),
            );
          },
        ),
      ),
    );
  }
}

getCurrentFlag(String currentLanguage) {
  if (currentLanguage == 'us') {
    currentLanguage = 'en';
  }
  for (Map flag in flags) {
    if (flag.containsValue(currentLanguage)) {
      return flag['image'];
    }
  }
}
