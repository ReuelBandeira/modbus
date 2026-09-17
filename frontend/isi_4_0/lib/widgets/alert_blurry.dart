import 'dart:ui';
import 'package:flutter/material.dart';

class BlurryDialog extends StatelessWidget {

  final String title;
  final String content;
  final String ButtonConfirmText;
  final String ButtonCancelText;
  final VoidCallback confirmCallBack;

  const BlurryDialog({
    super.key,
    required this.ButtonConfirmText,
    required this.ButtonCancelText,
    required this.title,
    required this.content,
    required this.confirmCallBack
  });

  final TextStyle textStyle = const TextStyle (color: Colors.black);

  @override
  Widget build(BuildContext context) {
    return BackdropFilter(
        filter: ImageFilter.blur(sigmaX: 6, sigmaY: 6),
        child:  AlertDialog(
          title: Text(title,style: textStyle,),
          content: Text(content, style: textStyle,),
          actions: <Widget>[
            TextButton(
              child: Text(ButtonConfirmText),
              onPressed: () {
                confirmCallBack();
                Navigator.of(context).pop();
              },
            ),
            TextButton(
              child: Text(ButtonCancelText),
              onPressed: () {
                Navigator.of(context).pop();
              },
            ),
          ],
        ));
  }
}