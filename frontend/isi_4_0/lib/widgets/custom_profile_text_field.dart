import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:isi_4_0/colors/color.dart';
import 'package:isi_4_0/utils/custom_error_filter_input.dart';
import 'package:isi_4_0/providers/lang.dart';
import 'package:provider/provider.dart';

class CustomProfileTextField extends StatefulWidget {
  final TextEditingController controller;
  final String labelText;
  final int? min;
  final int? max;
  final bool? isRequired;
  final bool? digitsOnly;
  final bool? ipOnly;
  final double? width;
  final bool? enabled;
  final bool? obscureText;
  final bool? isEmail;
  final bool? isUser;
  final Map<String, dynamic> translator;
  final String? error;
  final Function(String)? onChanged;
  const CustomProfileTextField({
    Key? key,
    required this.controller,
    required this.labelText,
    this.isEmail,
    this.isUser,
    this.min,
    this.max,
    this.isRequired,
    this.digitsOnly,
    this.ipOnly,
    this.width,
    this.enabled,
    this.obscureText,
    required this.translator,
    this.error,
    this.onChanged,
  }) : super(key: key);

  @override
  State createState() => _CustomProfileTextFieldState();
}

class _CustomProfileTextFieldState extends State<CustomProfileTextField> {
  bool showError = false;
  String errorString = '';
  bool currentEnabled = true;
  Map<String, dynamic> translator = {};
  String previousError = ' ';

  @override
  void initState() {
    setState(() {
      showError;
      errorString;
      translator;
      currentEnabled;
      previousError;
    });
    super.initState();
  }

  @override
  Widget build(BuildContext context) {
    double heightTextField = 50;
    double heightAlert = 28;

    if ((widget.error != null)) {
      if (widget.error != previousError) {
        previousError = widget.error!;
        showError = !(widget.error == '');
        errorString = widget.error!;
        setState(() {
          previousError;
          errorString;
          showError;
        });
      }
    }

    Map<String, dynamic> newTranslator = widget.translator;

    if ((newTranslator != translator) || (currentEnabled != widget.enabled)) {
      translator = newTranslator;
      currentEnabled = widget.enabled ?? true;
      setState(() {
        translator;
        currentEnabled;
      });
      CustomErrorFilterInput customFormatter = CustomErrorFilterInput(
          min: widget.min,
          max: widget.max,
          isRequired: widget.isRequired,
          digitsOnly: widget.digitsOnly,
          ipOnly: widget.ipOnly,
          enabled: widget.enabled,
          isEmail: widget.isEmail,
          isUser: widget.isUser,
          translator: translator,
          onError: (match, error) {
            setState(() {
              errorString = error;
              showError = match;
            });
          }
      );

      if (showError) {
        customFormatter.checkErrors(widget.controller.text);
      }
    }

    return SizedBox(
        height: showError ? heightTextField + heightAlert : heightTextField,
        width: widget.width,
        child: Column(
          children: [
            TextField(
              enabled: widget.enabled ?? true,
              obscureText: widget.obscureText ?? false,
              controller: widget.controller,
              // autovalidateMode: AutovalidateMode.onUserInteraction,
              style: const TextStyle(fontSize: 13),
              onChanged: widget.onChanged,
              inputFormatters: [
                  CustomErrorFilterInput(
                      min: widget.min,
                      max: widget.max,
                      isRequired: widget.isRequired,
                      digitsOnly: widget.digitsOnly,
                      ipOnly: widget.ipOnly,
                      enabled: widget.enabled,
                      isEmail: widget.isEmail,
                      isUser: widget.isUser,
                      translator: translator,
                      onError: (match, error) {
                        setState(() {
                          errorString = error;
                          showError = match;
                        });
                      }
                  ),
                // FilteringTextInputFormatter.deny(
                //   RegExp(r'\s'),
                // ),
                if (widget.digitsOnly != null && widget.digitsOnly == true)
                  FilteringTextInputFormatter.allow(RegExp(r'^-?[0-9]*'))
                else if (widget.ipOnly != null && widget.ipOnly == true)
                  FilteringTextInputFormatter.allow(RegExp(r'^[\d.]+')),
                // if (widget.isUser != null && widget.isUser!) FilteringTextInputFormatter.deny(RegExp(r'[ ]')),
                if (widget.isUser != null && widget.isUser!) LengthLimitingTextInputFormatter(50),
                if (widget.isEmail != null && widget.isEmail!) LengthLimitingTextInputFormatter(100),
              ],
              decoration: InputDecoration(
                filled: true,
                fillColor: CustomColors.whiteColorHigh,
                labelText: widget.labelText,
                labelStyle: MaterialStateTextStyle.resolveWith(
                      (Set<MaterialState> states) {
                    if (states.contains(MaterialState.disabled)) {
                      return const TextStyle(color: CustomColors.neutral400, fontSize: 13);
                    } else {
                      return const TextStyle(color: CustomColors.neutral700, fontSize: 13);
                    }
                  },
                ),
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
