import 'package:flutter/material.dart';
import 'package:isi_4_0/colors/color.dart';
import 'package:isi_4_0/utils/app_details.dart';

class FooterLogin extends StatelessWidget {
  final double height;
  final String version;
  const FooterLogin({super.key, required this.version, required this.height});

  @override
  Widget build(BuildContext context) {
    return Container(
      height: height,
      alignment: Alignment.centerRight,
      padding: const EdgeInsets.all(16),
      child: Text(
        version + AppDetails.appVertion.text,
        style: const TextStyle(
            fontSize: 10,
            fontWeight: FontWeight.normal,
            color: CustomColors.secondaryColorApp),
      ),
    );
  }
}
