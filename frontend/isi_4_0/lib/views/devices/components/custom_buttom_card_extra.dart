import 'package:flutter/material.dart';
import 'package:isi_4_0/colors/color.dart';

class CustomButtomCardExtra extends StatelessWidget {
  const CustomButtomCardExtra(
      {super.key,
      required this.height,
      required this.width,
      required this.field,
      required this.hover,
      required this.select,
      required this.showIcon,
      this.onPressedDelete});
  final double height;
  final double width;
  final String field;
  final bool hover;
  final bool select;
  final bool showIcon;
  final void Function()? onPressedDelete;
  @override
  Widget build(BuildContext context) {
    return Container(
      height: height,
      width: width,
      alignment: Alignment.centerLeft,
      padding: const EdgeInsets.only(bottom: 2),
      child: Row(
        children: [
          Card(
            elevation: 2,
            color: select
                ? CustomColors.primaryColorApp
                : hover
                    ? CustomColors.primary200
                    : CustomColors.whiteColorLow,
            shape: RoundedRectangleBorder(
              borderRadius: BorderRadius.circular(4.0),
            ),
            child: Container(
              height: height,
              width: showIcon ? width - 56 : width - 8,
              alignment: Alignment.centerLeft,
              padding: const EdgeInsets.only(left: 16),
              child: Text(
                field,
                style: TextStyle(
                    color: select
                        ? CustomColors.whiteColorHigh
                        : CustomColors.background600),
              ),
            ),
          ),
          Visibility(
            visible: showIcon,
            child: SizedBox(
              height: 48,
              width: 48,
              child: TextButton(
                onPressed: onPressedDelete,
                child: const Icon(
                  Icons.delete_forever_rounded,
                  color: CustomColors.greyDefault,
                ),
              ),
            ),
          )
        ],
      ),
    );
  }
}
