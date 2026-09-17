import 'package:flutter/material.dart';
import 'package:isi_4_0/colors/color.dart';
import 'package:isi_4_0/providers/lang.dart';
import 'package:isi_4_0/utils/app_details.dart';
import 'package:provider/provider.dart';

class CustomProfile extends StatelessWidget {
  final VoidCallback onTap;
  final String letter;
  final String? name;
  final String? role;
  final String? message;
  final String? imageName;
  const CustomProfile(
      {super.key,
      required this.onTap,
      required this.letter,
      this.name,
      this.role,
      this.message,
      this.imageName,
      Image? image});

  @override
  Widget build(BuildContext context) {
    return Row(
      children: [
        Column(
          mainAxisAlignment: MainAxisAlignment.center,
          crossAxisAlignment: CrossAxisAlignment.end,
          children: [
            Text(
              name!,
              style: const TextStyle(color: CustomColors.whiteColorLow, fontSize: 13),
            ),
            Padding(
              padding: const EdgeInsets.only(top: 2.0),
              child: Text(
                role!,
                style: const TextStyle(color: CustomColors.whiteColorLow, fontSize: 10),
              ),
            ),
          ],
        ),
        Tooltip(
          message: message ?? "",
          decoration: BoxDecoration(
            color: CustomColors.primaryColorApp.withOpacity(0.75),
            borderRadius: BorderRadius.circular(3),
          ),
          child: InkWell(
            onTap: onTap,
            child: Container(
              margin: const EdgeInsets.only(left: 8.0),
              padding: const EdgeInsets.all(8),
              alignment: Alignment.center,
              height: 32,
              width: 32,
              decoration: BoxDecoration(
                  borderRadius: BorderRadius.circular(90),
                  color: Colors.white70),
              child: Text(letter),
            ),
          ),
        ),
        Padding(
          padding: const EdgeInsets.only(right: 8.0),
          child: Consumer<LanguageProvider>(
            builder: (context, language, child) => InkWell(
              onTap: () async {
                showMenu(
                    context: context,
                    position: const RelativeRect.fromLTRB(60, 68, 0, 0),
                    items: [
                      PopupMenuItem(
                        child: Image.asset(
                          flags[0]['image'],
                          height: 20,
                          width: 20,
                        ),
                        onTap: () {
                          language.updateCurrentLanguage(flags[0]['lang']);
                        },
                      ),
                      PopupMenuItem(
                        child: Image.asset(
                          flags[1]['image'],
                          height: 20,
                          width: 20,
                        ),
                        onTap: () {
                          language.updateCurrentLanguage(flags[1]['lang']);
                        },
                      ),
                      PopupMenuItem(
                        child: Image.asset(
                          flags[2]['image'],
                          height: 20,
                          width: 20,
                        ),
                        onTap: () {
                          language.updateCurrentLanguage(flags[2]['lang']);
                        },
                      ),
                    ]);
              },
              child: Container(
                height: 36,
                width: 36,
                padding: const EdgeInsets.all(8),
                decoration: const BoxDecoration(
                  shape: BoxShape.circle,
                ),
                child: Image.asset(flagName(imageName)),
              ),
            ),
          ),
        )
      ],
    );
  }
}

String flagName(String? lang) {
  for (var flag in flags) {
    if (flag['lang'] == lang) return flag['image'];
  }
  return "flag_en.png";
}
