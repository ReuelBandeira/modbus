import 'package:flutter/material.dart';
import 'package:isi_4_0/colors/color.dart';

class CustomButtonVerticalHome extends StatelessWidget {
  final String toolTips;
  final List<String> buttonNames;
  final List<IconData> icons;
  final int index;
  final int currentPosition;
  final double heightButton;
  final bool minMenu;
  final VoidCallback onTap;

  const CustomButtonVerticalHome(
      {super.key,
      required this.toolTips,
      required this.buttonNames,
      required this.index,
      required this.icons,
      required this.currentPosition,
      required this.heightButton,
      required this.minMenu,
      required this.onTap});

  @override
  Widget build(BuildContext context) {
    return SizedBox(
      child: Container(
        margin: const EdgeInsets.all(8.0),
        height: heightButton,
        decoration: BoxDecoration(
            color: index == currentPosition
                ? CustomColors.blueColorLow.withOpacity(0.3)
                : Colors.transparent,
            borderRadius: BorderRadius.circular(4)),
        child: InkWell(
          onTap: onTap,
          child: Tooltip(
            message: minMenu ? "" : toolTips,
            decoration: BoxDecoration(
                color: CustomColors.primaryColorApp.withOpacity(0.75),
                borderRadius: BorderRadius.circular(3)),
            child: Row(
              children: [
                Padding(
                  padding: const EdgeInsets.only(left: 8.0),
                  child: Icon(
                    icons[index],
                    color: index == currentPosition
                        ? CustomColors.whiteColorLow.withOpacity(0.7)
                        : CustomColors.whiteColorLow,
                    size: 20,
                  ),
                ),
                Padding(
                  padding: const EdgeInsets.only(left: 8.0),
                  child: Visibility(
                    visible: minMenu,
                    child: Text(
                      buttonNames[index],
                      style: TextStyle(
                        color: index == currentPosition
                            ? CustomColors.whiteColorLow.withOpacity(0.7)
                            : CustomColors.whiteColorLow,
                      ),
                    ),
                  ),
                ),
              ],
            ),
          ),
        ),
      ),
    );
  }
}
