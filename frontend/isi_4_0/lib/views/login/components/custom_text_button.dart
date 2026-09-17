import 'package:flutter/material.dart';

class CustomTextButton extends StatelessWidget {
  final double height;
  final void Function()? onPressed;
  final String name;
  final Alignment alignment;
  const CustomTextButton(
      {super.key,
      required this.height,
      this.onPressed,
      required this.name,
      required this.alignment});

  @override
  Widget build(BuildContext context) {
    return Container(
      alignment: alignment,
      margin: const EdgeInsets.all(1),
      // padding: const EdgeInsets.only(left: 8),
      height: height,
      child: TextButton(
        onPressed: onPressed,
        child: Text(
          name,
          style: const TextStyle(fontSize: 12),
        ),
      ),
    );
  }
}
