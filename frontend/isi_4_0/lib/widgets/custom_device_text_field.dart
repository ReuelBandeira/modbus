import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:isi_4_0/colors/color.dart';
import 'package:isi_4_0/utils/custom_error_filter_input.dart';
import 'package:isi_4_0/providers/lang.dart';
import 'package:provider/provider.dart';

class CustomDeviceTextField extends StatefulWidget {
  final TextEditingController controller;
  final String labelText;
  final int? min;
  final int? max;
  final bool? isRequired;
  final bool? digitsOnly;
  final bool? ipOnly;
  const CustomDeviceTextField({
    Key? key,
    required this.controller,
    required this.labelText,
    this.min,
    this.max,
    this.isRequired,
    this.digitsOnly,
    this.ipOnly,
  }) : super(key: key);

  @override
  State createState() => _CustomDeviceTextFieldState();
}

class _CustomDeviceTextFieldState extends State<CustomDeviceTextField> {
  bool showError = false;
  String errorString = '';
  Map<String, dynamic> translator = {};

  @override
  void initState() {
    setState(() {
      showError;
      errorString;
      translator;
    });
    super.initState();
  }

  @override
  Widget build(BuildContext context) {
    double heightTextField = 50;
    double heightAlert = 14;

    Map<String, dynamic> newTranslator = (context).select(
            (LanguageProvider lang) => lang.getDataLanguage(lang.currentLanguage));

    if (newTranslator != translator) {
      translator = newTranslator;
      setState(() {
        translator;
      });
      CustomErrorFilterInput customFormatter = CustomErrorFilterInput(
          min: widget.min,
          max: widget.max,
          isRequired: widget.isRequired,
          digitsOnly: widget.digitsOnly,
          ipOnly: widget.ipOnly,
          translator: translator,
          onError: (match, error) {
            setState(() {
              errorString = error;
              showError = match;
            });
          }
      );

      customFormatter.checkErrors(widget.controller.text);
    }

    return SizedBox(
        height: showError ? heightTextField + heightAlert : heightTextField,
        child: Column(
          children: [
            TextField(
              controller: widget.controller,
              // autovalidateMode: AutovalidateMode.onUserInteraction,
              style: const TextStyle(fontSize: 13),
              inputFormatters: [
                  CustomErrorFilterInput(
                      min: widget.min,
                      max: widget.max,
                      isRequired: widget.isRequired,
                      digitsOnly: widget.digitsOnly,
                      ipOnly: widget.ipOnly,
                      translator: translator,
                      onError: (match, error) {
                        setState(() {
                          errorString = error;
                          showError = match;
                        });
                      }
                  ),
                FilteringTextInputFormatter.deny(
                  RegExp(r'\s'),
                ),
                if (widget.digitsOnly != null && widget.digitsOnly == true)
                  FilteringTextInputFormatter.allow(RegExp(r'^-?[0-9]*'))
                else if (widget.ipOnly != null && widget.ipOnly == true)
                  FilteringTextInputFormatter.allow(RegExp(r'^[\d.]+')),
              ],
              decoration: InputDecoration(
                filled: true,
                fillColor: CustomColors.whiteColorHigh,
                labelText: widget.labelText,
                labelStyle: const TextStyle(color: CustomColors.neutral700, fontSize: 13),
                border: OutlineInputBorder(
                    borderRadius: BorderRadius.circular(4),
                    borderSide: BorderSide.none),
                // errorText: validatePassword(controller.text),
              ),
            ),
            showError ? Container(
              height: heightAlert,
              padding: const EdgeInsets.only(top: 2, left: 8),
              alignment: Alignment.topLeft,
              child: Text(
                errorString,
                style: const TextStyle(fontSize: 10, color: CustomColors.error600),
              ),
            ) :
            Container()
          ],
        )
    );
  }
}
