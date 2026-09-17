import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:isi_4_0/colors/color.dart';
import 'package:isi_4_0/providers/lang.dart';
import 'package:isi_4_0/utils/custom_error_filter_input.dart';
import 'package:isi_4_0/utils/enum_login.dart';
import 'package:provider/provider.dart';

class CustomTextFormField extends StatefulWidget {
  final TextEditingController controller;
  final void Function(String)? onFieldSubmitted;
  final String? Function(String?)? validator;
  final CurrentPage currentPage;
  final CurrentTextField currentTextField;
  final String? hintText;
  final Widget? prefixIcon;
  final bool? obscureText;
  final List<TextInputFormatter>? inputFormatters;
  final String? keyText;
  final String? error;
  final Function(String)? onChanged;

  const CustomTextFormField({
    super.key,
    required this.controller,
    this.onFieldSubmitted,
    required this.currentPage,
    this.hintText,
    this.prefixIcon,
    this.obscureText,
    required this.currentTextField,
    this.validator,
    this.inputFormatters,
    this.keyText,
    this.error,
    this.onChanged,
  });

  @override
  State createState() => _CustomTextFormFieldState();
}

class _CustomTextFormFieldState extends State<CustomTextFormField> {
  bool showError = false;
  String errorString = '';
  Map<String, dynamic> translator = {};
  bool isRequired = false;
  bool isEmail = false;
  bool isUser = false;
  String previousError = ' ';

  @override
  void initState() {
    if ((widget.currentTextField == CurrentTextField.user) ||
        (widget.currentTextField == CurrentTextField.email) ||
        (widget.currentTextField == CurrentTextField.password) ||
        (widget.currentTextField == CurrentTextField.confirmPassword) ||
        (widget.currentTextField == CurrentTextField.addressIP) ||
        (widget.currentTextField == CurrentTextField.port)) {
      if (widget.currentPage != CurrentPage.settings) {
        isRequired = true;
      } else {
        if ((widget.currentTextField == CurrentTextField.addressIP) ||
            (widget.currentTextField == CurrentTextField.port)) {
          isRequired = true;
        }
      }
    }
    if (widget.currentTextField == CurrentTextField.email) {
      isEmail = true;
    }
    if (widget.currentTextField == CurrentTextField.user) {
      if (widget.currentPage != CurrentPage.settings) {
        isUser = true;
      }
    }
    setState(() {
      showError;
      errorString;
      translator;
      isRequired;
      isEmail;
      isUser;
      previousError;
    });
    super.initState();
  }

  @override
  Widget build(BuildContext context) {
    double heightTextField = 65;
    double heightAlert = 14;

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

    Map<String, dynamic> newTranslator = (context).select(
            (LanguageProvider lang) => lang.getDataLanguage(lang.currentLanguage));

    if (newTranslator != translator) {
      translator = newTranslator;
      setState(() {
        translator;
      });
      CustomErrorFilterInput customFormatter = CustomErrorFilterInput(
          isRequired: isRequired,
          isEmail: isEmail,
          isUser: isUser,
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
        child: Padding(
          padding: const EdgeInsets.all(8.0),
          child: Column(
            children: [
              Container(
                padding: widget.prefixIcon != null
                    ? const EdgeInsets.only(left: 0)
                    : const EdgeInsets.only(left: 16),
                decoration: BoxDecoration(
                    color: CustomColors.whiteColorHigh, borderRadius: BorderRadius.circular(6)
                ),
                child: Semantics(
                  label: widget.keyText,
                  child: TextFormField(
                    obscureText: widget.obscureText ?? false,
                    controller: widget.controller,
                    onFieldSubmitted: widget.onFieldSubmitted,
                    validator: widget.validator,
                    onChanged: widget.onChanged,
                    inputFormatters: [
                      CustomErrorFilterInput(
                          isRequired: isRequired,
                          isEmail: isEmail,
                          isUser: isUser,
                          translator: translator,
                          onError: (match, error) {
                            setState(() {
                              errorString = error;
                              showError = match;
                            });
                          }
                      ),
                      // if (isUser) FilteringTextInputFormatter.deny(RegExp(r'[ ]')),
                      if (isUser) LengthLimitingTextInputFormatter(50),
                      if (isEmail) LengthLimitingTextInputFormatter(100),
                      ...?widget.inputFormatters
                    ],
                    decoration: InputDecoration(
                      prefixIcon: widget.prefixIcon,
                      hintText: widget.hintText,
                      contentPadding: (widget.currentPage != CurrentPage.settings) ? const EdgeInsets.only(top: 14) : null,
                      border: InputBorder.none,
                      errorStyle: const TextStyle(
                        height: 0,
                        fontSize: 0.0,
                      ),
                    ),
                  )
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
          ),
        ),
    );
  }
}
