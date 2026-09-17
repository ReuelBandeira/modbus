import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:isi_4_0/colors/color.dart';

class CustomTextField extends StatelessWidget {
  final String? hint;
  final double width;
  final double height;
  final Color backgroundColor;
  final Color textColor;
  final bool showIcon;
  final VoidCallback onPressedIcon;
  final Widget icon;
  final double splashRadiusIcon;
  final TextEditingController controller;
  final bool obscureText;
  final double borderRadius;
  final FocusNode? focusNode;
  final bool? enable;
  final void Function(String)? onSubmitted;
  final void Function(String)? onChanged;
  final List<TextInputFormatter>? inputFormatters;
  const CustomTextField({
    Key? key,
    Color? iconColor,
    required this.icon,
    required this.onPressedIcon,
    required this.hint,
    required this.width,
    required this.height,
    required this.backgroundColor,
    required this.textColor,
    required this.showIcon,
    required this.splashRadiusIcon,
    required this.controller,
    required this.borderRadius,
    required this.obscureText,
    this.focusNode,
    this.inputFormatters,
    this.enable,
    this.onSubmitted,
    this.onChanged,
  }) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return Container(
      width: width,
      height: height + 10,
      decoration: BoxDecoration(
        color: CustomColors.whiteColorHigh,
        borderRadius: BorderRadius.circular(borderRadius),
      ),
      child: TextField(
        enabled: enable,
        textInputAction: TextInputAction.next,
        obscureText: obscureText,
        controller: controller,
        cursorColor: CustomColors.greyColorHigh,
        inputFormatters: inputFormatters,
        focusNode: focusNode,
        onSubmitted: onSubmitted,
        onChanged: onChanged,
        style: const TextStyle(fontSize: 14, color: CustomColors.greyColorHigh),
        decoration: showIcon
            ? InputDecoration(
                focusColor: CustomColors.primaryColorApp,
                prefixIcon: IconButton(
                  splashRadius: splashRadiusIcon,
                  icon: icon,
                  color: CustomColors.greyColorHigh,
                  onPressed: onPressedIcon,
                ),
                border: InputBorder.none,
                hintText: hint,
                hintStyle: TextStyle(color: textColor, fontSize: 12),
              )
            : InputDecoration(
                hintStyle: TextStyle(color: textColor),
                suffixIcon: const Icon(
                  Icons.search,
                  color: CustomColors.primaryColorApp,
                ),
                border: InputBorder.none,
                enabledBorder: const UnderlineInputBorder(
                    borderSide: BorderSide(color: CustomColors.primaryColorApp)),
                focusedBorder: const UnderlineInputBorder(
                  borderSide: BorderSide(color: CustomColors.primaryColorApp),
                ),
                hintText: hint,
              ),
      ),
    );
  }
}
