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

class SignIn extends StatefulWidget {
  const SignIn({super.key});

  @override
  State<SignIn> createState() => _SignInState();
}

class _SignInState extends State<SignIn> {
  AuthService authService = AuthService.instance;

  final formKeySignIn = GlobalKey<FormState>();

  final TextEditingController emailController = TextEditingController();
  final TextEditingController passwordController = TextEditingController();

  @override
  void initState() {
    initialSettings();
    super.initState();
  }

  Future<void> initialSettings() async {
    await authService.refreshPage('/');
  }

  @override
  void dispose() {
    emailController.dispose();
    passwordController.dispose();
    super.dispose();
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
              rightSide(mainHeight, mainWidth, formKeySignIn, emailController,
                  passwordController, languageProvider, loginProvider),
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
  LanguageProvider languageProvider,
  LoginViewModel loginProvider,
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
                body(formKey, emailController, passwordController, words,
                    sizeRight, loginProvider, context, languageProvider),
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
    words,
    BoxConstraints sizeRight,
    LoginViewModel loginProvider,
    BuildContext context,
    LanguageProvider languageProvider) {
  return SizedBox(
    height: sizeRight.maxHeight > 600 ? sizeRight.maxHeight - 96 : 600,
    child: Column(
      crossAxisAlignment: CrossAxisAlignment.center,
      mainAxisAlignment: MainAxisAlignment.center,
      children: [
        customLog(),
        customForm(formKey, emailController, passwordController, words,
            sizeRight, loginProvider, context, languageProvider),
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
    words,
    BoxConstraints sizeRight,
    LoginViewModel loginProvider,
    BuildContext context,
    LanguageProvider languageProvider) {
  bool iconPasswordState = true;
  bool iconPasswordHover = false;

  return SizedBox(
    width: sizeRight.maxWidth * 0.5,
    child: Form(
      key: formKey,
      child: Column(
        children: [
          const SizedBox(height: 40),
          CustomTextFormField(
              controller: emailController,
              currentTextField: CurrentTextField.email,
              currentPage: CurrentPage.signIn,
              hintText: words["email"],
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
              validator: (value) {
                return MessageValidator.instance.getMessage(
                    CurrentPage.signIn, value, CurrentTextField.email, words);
              }),
          StatefulBuilder(
            builder: (context, passwordState) {
              return CustomTextFormField(
                  controller: passwordController,
                  currentTextField: CurrentTextField.password,
                  currentPage: CurrentPage.signIn,
                  hintText: words["password"],
                  obscureText: iconPasswordState,
                  inputFormatters: [
                    FilteringTextInputFormatter.deny(RegExp(r'\s'))
                  ],
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
                    submitForm(formKey, emailController, passwordController,
                        loginProvider, context, languageProvider);
                  },
                  validator: (value) {
                    return MessageValidator.instance.getMessage(
                        CurrentPage.signIn,
                        value,
                        CurrentTextField.password,
                        words);
                  });
            },
          ),
          CustomTextButton(
            height: 36,
            alignment: Alignment.centerLeft,
            name: words['forgotPassword'],
            onPressed: () {
              loginProvider.recoverView(context);
            },
          ),
          const SizedBox(height: 40),
          customButtom(words['enter'], formKey, emailController,
              passwordController, loginProvider, context, languageProvider),
          const SizedBox(height: 20),
          CustomTextButton(
            height: 36,
            alignment: Alignment.center,
            name: words['register'],
            onPressed: () {
              loginProvider.registerView(context);
            },
          ),
        ],
      ),
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
    int code = await authService.readUser(
        emailController.text, passwordController.text);
    if (code == 200) {
      // ignore: use_build_context_synchronously
      await loginProvider.homeView(context);
    } else {
      if (context.mounted) {
        alertMessage.showError(context, messageAPI[code.toString()]);
      }
    }
  }
}
