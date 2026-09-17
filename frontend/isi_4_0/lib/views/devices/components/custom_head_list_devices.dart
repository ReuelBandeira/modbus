import 'package:flutter/material.dart';
import 'package:isi_4_0/colors/color.dart';
import 'package:isi_4_0/providers/lang.dart';

class CustomHeadListDevices extends StatelessWidget {
  const CustomHeadListDevices({super.key, required this.language});
  final LanguageProvider language;
  @override
  Widget build(BuildContext context) {
    Map<String, dynamic> lang =
        language.getDataLanguage(language.currentLanguage)['headDeviceList'];

    return Padding(
      padding: const EdgeInsets.only(top: 24.0, left: 4),
      child: SingleChildScrollView(
        scrollDirection: Axis.horizontal,
        child: Row(
          children: lang.entries
              .map((field) => Container(
                  height: 48,
                  width: 180,
                  alignment: Alignment.centerLeft,
                  padding: const EdgeInsets.only(left: 8),
                  decoration: const BoxDecoration(
                      color: CustomColors.lineColors,
                      border: Border(
                          bottom: BorderSide(color: Colors.black54, width: 4))
                  ),
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
