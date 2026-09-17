import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:isi_4_0/colors/color.dart';
import 'package:isi_4_0/providers/lang.dart';
import 'package:isi_4_0/repository/db_api.dart';
import 'package:isi_4_0/utils/alert_message.dart';
import 'package:isi_4_0/utils/enum_login.dart';
import 'package:isi_4_0/utils/message_validator.dart';
import 'package:isi_4_0/utils/text_form_field.dart';
import 'package:isi_4_0/viewmodel/login_view_model.dart';
import 'package:isi_4_0/views/login/components/custom_text_button.dart';
import 'package:isi_4_0/views/login/components/footer_login.dart';
import 'package:isi_4_0/views/login/components/header_login.dart';
import 'package:isi_4_0/widgets/custom_rounded_button.dart';
import 'package:provider/provider.dart';

class Recover extends StatefulWidget {
  const Recover({super.key});

  @override
  State<Recover> createState() => _RecoverState();
}

class _RecoverState extends State<Recover> {
  AuthService authService = AuthService.instance;

  final formKeyRecover = GlobalKey<FormState>();

  final TextEditingController emailController = TextEditingController();
  final TextEditingController passwordController = TextEditingController();
  final TextEditingController confirmPasswordController =
      TextEditingController();

  String errorMatch = '';

  @override
  void initState() {
    initialSettings();
    setState(() {
      errorMatch;
    });
    super.initState();
  }

  Future<void> initialSettings() async {
    await authService.refreshPage('/recover');
  }

  @override
  void dispose() {
    emailController.dispose();
    passwordController.dispose();
    confirmPasswordController.dispose();
    super.dispose();
  }

  void setError(String error) {
    errorMatch = error;
    setState(() {
      errorMatch;
    });
  }

  @override
  Widget build(BuildContext context) {
    final double mainHeight = MediaQuery.of(context).size.height;
    final double mainWidth = MediaQuery.of(context).size.width;

    return Scaffold(
      body: Consumer2<LoginViewModel, LanguageProvider>(
        builder: (context, loginProvider, languageProvider, child) {
          return Row(
            children: [
              Visibility(
                  visible: (mainWidth - mainWidth * 0.6) > 300,
                  child: leftSide(mainHeight, mainWidth)),
              rightSide(
                  mainHeight,
                  mainWidth,
                  formKeyRecover,
                  emailController,
                  passwordController,
                  confirmPasswordController,
                  languageProvider,
                  loginProvider,
                  errorMatch,
                  setError),
            ],
          );
        },
      ),
    );
  }
}

Widget leftSide(double mainHeight, double mainWidth) {
  return SizedBox(
    height: mainHeight,
    width: mainWidth * 0.6,
    child: LayoutBuilder(
      builder: (context, sizeLeft) => Container(
        height: sizeLeft.maxHeight,
        width: sizeLeft.maxWidth,
        decoration: const BoxDecoration(
          image: DecorationImage(
            image: AssetImage("assets/bg_login.png"),
            fit: BoxFit.fill,
          ),
        ),
      ),
    ),
  );
}

Widget rightSide(
  double mainHeight,
  double mainWidth,
  GlobalKey<FormState> formKey,
  TextEditingController emailController,
  TextEditingController passwordController,
  TextEditingController confirmPasswordController,
  LanguageProvider languageProvider,
  LoginViewModel loginProvider,
  String error,
  Function(String) setError,
) {
  return SizedBox(
    height: mainHeight,
    width: (mainWidth - mainWidth * 0.6) > 300 ? mainWidth * 0.4 : mainWidth,
    child: LayoutBuilder(
      builder: (context, sizeRight) {
        var words = languageProvider
            .getDataLanguage(languageProvider.currentLanguage)["dictionary"];

        return SizedBox(
          height: sizeRight.maxHeight,
          child: SingleChildScrollView(
            child: Column(
              mainAxisAlignment: MainAxisAlignment.spaceBetween,
              children: [
                header(48, sizeRight, languageProvider),
                body(
                    formKey,
                    emailController,
                    passwordController,
                    confirmPasswordController,
                    words,
                    sizeRight,
                    loginProvider,
                    context,
                    languageProvider,
                    error,
                    setError),
                footer(words["version"])
              ],
            ),
          ),
        );
      },
    ),
  );
}

Widget header(double height, BoxConstraints sizeRight,
    LanguageProvider languageProvider) {
  return HeaderLogin(
      height: height,
      width: sizeRight.maxWidth,
      languageProvider: languageProvider);
}

Widget body(
    GlobalKey<FormState> formKey,
    TextEditingController emailController,
    TextEditingController passwordController,
    TextEditingController confirmPasswordController,
    words,
    BoxConstraints sizeRight,
    LoginViewModel loginProvider,
    BuildContext context,
    LanguageProvider languageProvider,
    String error,
    Function(String) setError) {
  return SizedBox(
    height: sizeRight.maxHeight > 666 ? sizeRight.maxHeight - 96 : 666,
    child: Column(
      crossAxisAlignment: CrossAxisAlignment.center,
      mainAxisAlignment: MainAxisAlignment.center,
      children: [
        customLog(),
        customForm(
            formKey,
            emailController,
            passwordController,
            confirmPasswordController,
            words,
            sizeRight,
            loginProvider,
            context,
            languageProvider,
            error,
            setError),
      ],
    ),
  );
}

Widget customLog() {
  return Container(
    height: 150,
    width: 150,
    alignment: Alignment.topCenter,
    decoration: const BoxDecoration(
      image: DecorationImage(
        image: AssetImage("assets/logo.png"),
      ),
    ),
  );
}

Widget customForm(
    GlobalKey<FormState> formKey,
    TextEditingController emailController,
    TextEditingController passwordController,
    TextEditingController confirmPasswordController,
    words,
    BoxConstraints sizeRight,
    LoginViewModel loginProvider,
    BuildContext context,
    LanguageProvider languageProvider,
    String error,
    Function(String) setError) {
  bool iconPasswordState = true;
  bool iconPasswordHover = false;
  // bool iconConfirmPasswordHover = false;
  AlertMessage alertMessage = AlertMessage();

  List<String> passwords = List.filled(2, '');

  return SizedBox(
    width: sizeRight.maxWidth * 0.5,
    child: Form(
      key: formKey,
      child: StatefulBuilder(builder: (context, passwordState) {
        return Column(
          children: [
            const SizedBox(height: 40),
            CustomTextFormField(
              controller: emailController,
              currentTextField: CurrentTextField.email,
              currentPage: CurrentPage.recover,
              hintText: words["email"],
              validator: (value) {
                return MessageValidator.instance.getMessage(
                    CurrentPage.recover, value, CurrentTextField.email, words);
              },
              inputFormatters: [
                FilteringTextInputFormatter.deny(RegExp(r'\s'))
              ],
              prefixIcon: GestureDetector(
                child: const Icon(
                  Icons.email_outlined,
                  size: 18,
                  color: CustomColors.greyColorHigh,
                ),
              ),
              onFieldSubmitted: (value) {
                submitForm(formKey, emailController, passwordController,
                    loginProvider, context, languageProvider);
              },
            ),
            const SizedBox(height: 20),
            CustomTextFormField(
              controller: passwordController,
              currentTextField: CurrentTextField.password,
              currentPage: CurrentPage.recover,
              obscureText: iconPasswordState,
              hintText: words["password"],
              inputFormatters: [
                FilteringTextInputFormatter.deny(RegExp(r'\s'))
              ],
              validator: (value) {
                if (passwordController.text.isNotEmpty &&
                    confirmPasswordController.text.isNotEmpty) {
                  if (passwordController.text !=
                      confirmPasswordController.text) {
                    setError(words['msg_passwords_not_match']);
                    if (context.mounted) {
                      alertMessage.showError(context, words['msg_passwords_not_match']);
                    }
                    return '';
                  } else {
                    return null;
                  }
                } else {
                  return MessageValidator.instance.getMessage(
                      CurrentPage.recover,
                      value,
                      CurrentTextField.password,
                      words);
                }
              },
              onChanged: (text) {
                setError('');
              },
              prefixIcon: MouseRegion(
                child: GestureDetector(
                  child: Container(
                    height: 20,
                    width: 20,
                    margin: const EdgeInsets.all(4),
                    decoration: BoxDecoration(
                        borderRadius: BorderRadius.circular(90),
                        color: iconPasswordHover ? CustomColors.hoverColorLogin : null),
                    child: Icon(
                      iconPasswordState
                          ? Icons.visibility_off_outlined
                          : Icons.remove_red_eye_outlined,
                      size: 20,
                      color: CustomColors.greyColorHigh,
                    ),
                  ),
                  onTap: () {
                    passwordState(() {
                      iconPasswordState = !iconPasswordState;
                    });
                  },
                ),
                onEnter: (event) {
                  passwordState(() {
                    iconPasswordHover = true;
                  });
                },
                onExit: (event) {
                  passwordState(() {
                    iconPasswordHover = false;
                  });
                },
              ),
              onFieldSubmitted: (value) {
                passwordState(
                  () {
                    passwords.first = passwordController.text;
                    passwords.last = confirmPasswordController.text;
                  },
                );

                submitForm(formKey, emailController, passwordController,
                    loginProvider, context, languageProvider);
              },
            ),
            CustomTextFormField(
              controller: confirmPasswordController,
              currentTextField: CurrentTextField.confirmPassword,
              currentPage: CurrentPage.recover,
              obscureText: iconPasswordState,
              hintText: words["confirmPassword"],
              inputFormatters: [
                FilteringTextInputFormatter.deny(RegExp(r'\s'))
              ],
              validator: (value) {
                if (passwordController.text.isNotEmpty &&
                    confirmPasswordController.text.isNotEmpty) {
                  if (passwordController.text !=
                      confirmPasswordController.text) {
                    return '';//words['msg_passwords_not_match'];
                  } else {
                    return null;
                  }
                } else {
                  return MessageValidator.instance.getMessage(
                      CurrentPage.recover,
                      value,
                      CurrentTextField.confirmPassword,
                      words);
                }
              },
              error: error,
              onChanged: (text) {
                setError('');
              },
              prefixIcon: Container(
                height: 20,
                width: 20,
                margin: const EdgeInsets.all(4),
              ),
            ),
            const SizedBox(height: 40),
            customButtom(words['recover'], formKey, emailController,
                passwordController, loginProvider, context, languageProvider),
            const SizedBox(height: 20),
            CustomTextButton(
              height: 36,
              alignment: Alignment.center,
              name: words['goBack'],
              onPressed: () {
                loginProvider.signInView(context);
              },
            ),
          ],
        );
      }),
    ),
  );
}

Widget customButtom(
    word,
    GlobalKey<FormState> formKey,
    TextEditingController emailController,
    TextEditingController passwordController,
    LoginViewModel loginProvider,
    BuildContext context,
    LanguageProvider languageProvider) {
  return CustomRoundedButton(
    textName: word,
    height: 32,
    width: 100,
    fontSize: 12,
    isSelected: true,
    textColorActived: CustomColors.whiteColorLow,
    textColorInactive: CustomColors.whiteColorLow,
    splashColor: CustomColors.primaryColorApp,
    backgroundColorActived: CustomColors.primaryColorApp,
    backgroundColorInactive: CustomColors.primaryColorApp,
    borderRadiusValue: 30,
    onTap: () {
      submitForm(formKey, emailController, passwordController, loginProvider,
          context, languageProvider);
    },
  );
}

Widget footer(version) {
  return FooterLogin(version: version, height: 48);
}

Future<void> submitForm(
    GlobalKey<FormState> formKey,
    TextEditingController emailController,
    TextEditingController passwordController,
    LoginViewModel loginProvider,
    BuildContext context,
    LanguageProvider languageProvider) async {
  if (formKey.currentState!.validate()) {
    AuthService authService = AuthService.instance;
    AlertMessage alertMessage = AlertMessage();
    var messageAPI = languageProvider
        .getDataLanguage(languageProvider.currentLanguage)["api"];
    int code = await authService.updateUser(
        emailController.text, passwordController.text);
    if (code == 200) {
      if (context.mounted) {
        alertMessage.showSuccess(context, messageAPI["recover_200"]);
      }
      // ignore: use_build_context_synchronously
      await loginProvider.signInView(context);
    } else {
      if (context.mounted) {
        alertMessage.showError(context, messageAPI[code.toString()]);
      }
    }
  }
}
