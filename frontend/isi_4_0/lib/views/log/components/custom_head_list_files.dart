import 'package:flutter/material.dart';
import 'package:isi_4_0/colors/color.dart';
import 'package:isi_4_0/providers/lang.dart';

class CustomHeadListFiles extends StatelessWidget {
  const CustomHeadListFiles({super.key, required this.language});
  final LanguageProvider language;
  @override
  Widget build(BuildContext context) {
    Map<String, dynamic> lang =
        language.getDataLanguage(language.currentLanguage)['headFileList'];

    return Padding(
      padding: const EdgeInsets.only(top: 24.0),
      child: SingleChildScrollView(
        scrollDirection: Axis.horizontal,
        child: Row(
          mainAxisAlignment: MainAxisAlignment.spaceBetween,
          children: lang.entries
              .map((field) => Container(
                  height: 48,
                  width: field.key == 'log' ? 780 : 300,
                  alignment: Alignment.centerLeft,
                  padding: const EdgeInsets.only(left: 8),
                  decoration: const BoxDecoration(
                      color: CustomColors.lineColors,
                      border: Border(
                          bottom: BorderSide(color: Colors.black54, width: 4))),
                  child: Text(field.value,
                      style: const TextStyle(
                          fontWeight: FontWeight.bold,
                          color: CustomColors.primaryColorApp))))
              .toList(),
        ),
      ),
    );
  }
}
