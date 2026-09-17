import 'package:flutter/material.dart';
import 'package:isi_4_0/colors/color.dart';

class CustomRowsFiles extends StatelessWidget {
  const CustomRowsFiles(
      {super.key,
      required this.nameFile,
      this.onPressed,
      this.tooltipsMessage,
      required this.dateFile});
  final String nameFile;
  final String dateFile;
  final String? tooltipsMessage;
  final void Function()? onPressed;
  @override
  Widget build(BuildContext context) {
    const double height = 40;
    const double width = 1080;
    return Container(
      height: height,
      margin: const EdgeInsets.only(top: 2.0),
      child: SingleChildScrollView(
        scrollDirection: Axis.horizontal,
        child: Row(
          children: [
            Container(
              height: height,
              width: width - 300,
              alignment: Alignment.centerLeft,
              padding: const EdgeInsets.only(left: 10.0),
              decoration: const BoxDecoration(
                  border: Border(
                      bottom: BorderSide(
                          color: CustomColors.lineColors, width: 2))),
              child: Text(nameFile),
            ),
            Container(
              height: height,
              width: 252,
              alignment: Alignment.centerLeft,
              padding: const EdgeInsets.only(left: 12.0),
              decoration: const BoxDecoration(
                  border: Border(
                      bottom: BorderSide(
                          color: CustomColors.lineColors, width: 2))),
              child: Text(dateFile),
            ),
            Tooltip(
              message: tooltipsMessage ?? "",
              decoration: BoxDecoration(
                color: CustomColors.primaryColorApp.withOpacity(0.75),
                borderRadius: BorderRadius.circular(90),
              ),
              child: Container(
                height: height,
                width: height,
                decoration: const BoxDecoration(
                    border: Border(
                        bottom: BorderSide(
                            color: CustomColors.lineColors, width: 2))),
                child: TextButton(
                  style: ButtonStyle(
                      padding:
                          MaterialStateProperty.all(const EdgeInsets.all(0)),
                      shape: MaterialStateProperty.all(RoundedRectangleBorder(
                          borderRadius: BorderRadius.circular(90)))),
                  onPressed: onPressed,
                  child: const Icon(
                    Icons.file_download_outlined,
                    size: 20,
                  ),
                ),
              ),
            )
          ],
        ),
      ),
    );
  }
}
