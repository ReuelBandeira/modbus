import 'package:flutter/material.dart';
import 'package:isi_4_0/colors/color.dart';

class DropbuttonVerticalWithButton extends StatelessWidget {
  final double? heightContainer;
  final double? widthContainer;
  final double? widthDropdownButton;
  final double? widthButton;
  final double? paddingLeft;
  final Color? backgroundColor;
  final Function()? onPressedIcon;
  final List<String> dropList;
  final String hintDrop;
  final String? tooltipsMessage;
  final Function(String?)? onChanged;

  const DropbuttonVerticalWithButton({
    Key? key,
    this.heightContainer,
    this.widthContainer,
    this.widthDropdownButton,
    this.widthButton,
    this.backgroundColor,
    this.paddingLeft,
    this.tooltipsMessage,
    required this.onPressedIcon,
    required this.dropList,
    required this.hintDrop,
    required this.onChanged,
  }) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return Container(
      height: heightContainer,
      width: widthContainer,
      color: backgroundColor,
      child: Row(
        children: [
          Container(
            height: heightContainer,
            width: widthDropdownButton ?? 0,
            padding: EdgeInsets.only(
                left: paddingLeft ?? 0, top: 0, right: 0, bottom: 0),
            child: DropdownButton(
              items: dropList.map(
                (String topic) {
                  return DropdownMenuItem(
                    value: topic,
                    child: Text(
                      topic,
                      style: const TextStyle(fontSize: 13),
                    ),
                  );
                },
              ).toList(),
              onChanged: onChanged,
              underline: const SizedBox.shrink(),
              hint: Text(
                hintDrop,
                style: const TextStyle(color: Colors.black54, fontSize: 13),
              ),
              isExpanded: true,
            ),
          ),
          Tooltip(
            message: tooltipsMessage ?? "",
            decoration: BoxDecoration(
              color: CustomColors.primaryColorApp.withOpacity(0.75),
              borderRadius: BorderRadius.circular(3),
            ),
            child: Container(
              height: heightContainer,
              width: widthButton ?? 0,
              color: backgroundColor,
              child: TextButton(
                onPressed: onPressedIcon,
                child: const Icon(
                  Icons.more_vert,
                  size: 20,
                ),
              ),
            ),
          )
        ],
      ),
    );
  }
}
