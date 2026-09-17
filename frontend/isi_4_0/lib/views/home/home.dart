import 'dart:async';
import 'dart:convert' as convert;

import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:isi_4_0/colors/color.dart';
import 'package:isi_4_0/providers/lang.dart';
import 'package:isi_4_0/providers/profile.dart';
import 'package:isi_4_0/repository/db_api.dart';
import 'package:isi_4_0/utils/alert_message.dart';
import 'package:isi_4_0/utils/navigator_service.dart';
import 'package:isi_4_0/viewmodel/device_view_model.dart';
import 'package:isi_4_0/viewmodel/home_view_model.dart';
import 'package:isi_4_0/viewmodel/log_view_model.dart';
import 'package:isi_4_0/viewmodel/settings_view_model.dart';
import 'package:isi_4_0/views/devices/pageview_devices.dart';
import 'package:isi_4_0/views/log/log.dart';
import 'package:isi_4_0/views/settings/settings.dart';
import 'package:isi_4_0/widgets/custom_button_vertical_home.dart';
import 'package:isi_4_0/widgets/custom_profile.dart';
import 'package:isi_4_0/widgets/custom_profile_text_field.dart';
import 'package:isi_4_0/widgets/custom_rounded_button.dart';
import 'package:provider/provider.dart';

class Home extends StatefulWidget {
  const Home({Key? key}) : super(key: key);
  @override
  State<Home> createState() => _HomeState();
}

class _HomeState extends State<Home> {
  AuthService authService = AuthService.instance;
  List<String> profile = List<String>.filled(3, "");
  String letter = "";

  late BuildContext te;

  var buttonIcons = [
    Icons.devices,
    Icons.feed,
    Icons.settings_outlined,
  ];

  bool extend = false;
  bool shortMenu = false;
  var currentPosition = 0;
  late IconData buttonMenu = Icons.menu;
  final pageController = PageController();
  final double heightButton = 32;
  final double widthButtonMin = 52;
  final double widthButtonMax = 200;
  final double widthMin = 501;
  final double heightHeader = 60;
  final double buttonHeight = 65;

  @override
  void initState() {
    super.initState();
    initialSettings();
  }

  Future<void> initialSettings() async {
    await authService.refreshPage('/home');
    String login = await authService.getLoginDetails();
    Map<String, dynamic> data = convert.jsonDecode(login);
    setState(() {
      profile[0] = data["name"];
      profile[1] = data["email"];
      profile[2] = data["id"].toString();
      letter = profile[0][0];
    });
  }

  @override
  Widget build(BuildContext context) {
    Map<String, dynamic> translator = (context).select(
        (LanguageProvider lang) => lang.getDataLanguage(lang.currentLanguage));

    List<String> toolTips = [
      translator["tooltips"]["extendMenu"],
      translator["tooltips"]["collapseMenu"],
      translator["tooltips"]["devices"],
      translator["tooltips"]["log"],
      translator["tooltips"]["settings"],
    ];

    Future<void> actionMenu(
        int select,
        List<String> profile,
        Map<String, dynamic> translator,
        HomeViewModel homeProvider,
        String currentLanguage,
        AuthService authService) async {
      switch (select) {
        case 0:
          bool result = await editProfile(context, profile, translator, authService);
          if (result) {
            initialSettings();
          }
          break;
        case 1:
          homeProvider.updateLanguage(profile[2], currentLanguage);
          NavigationService().logOut();
          // homeProvider.loginView(context);
          break;
        default:
      }
    }

    return NotificationListener<SizeChangedLayoutNotification>(
      onNotification: (notification) {
        return true;
      },
      child: SizeChangedLayoutNotifier(
        child: LayoutBuilder(
          builder: (context, constraints) {
            return Consumer2<LanguageProvider, HomeViewModel>(
              builder: (context, language, homeProvider, child) {
                {
                  return Scaffold(
                    backgroundColor: CustomColors.whiteColorLow,
                    body: OverflowBox(
                      minHeight: 0,
                      minWidth: 0,
                      maxHeight: constraints.maxHeight,
                      maxWidth: constraints.maxWidth,
                      child: Column(
                        children: [
                          Card(
                            color: CustomColors.primaryColorApp,
                            elevation: 4,
                            margin: const EdgeInsets.all(0),
                            shadowColor: CustomColors.secondaryColorApp,
                            child: Container(
                              height: heightHeader,
                              width: constraints.maxWidth,
                              decoration: const BoxDecoration(
                                color: CustomColors.primaryColorApp,
                                border: Border(
                                    bottom: BorderSide(
                                        color: Colors.black12, width: 2)),
                              ),
                              child: Row(
                                mainAxisAlignment:
                                    MainAxisAlignment.spaceBetween,
                                children: [
                                  Row(
                                    mainAxisAlignment: MainAxisAlignment.center,
                                    children: [
                                      Visibility(
                                        visible: extend,
                                        child: SizedBox(
                                          height: heightHeader,
                                          width: 144,
                                        ),
                                      ),
                                      InkWell(
                                        child: Tooltip(
                                          message: extend
                                              ? toolTips[1]
                                              : toolTips[0],
                                          decoration: BoxDecoration(
                                            color: CustomColors.primaryColorApp
                                                .withOpacity(0.75),
                                            borderRadius:
                                                BorderRadius.circular(3),
                                          ),
                                          child: Padding(
                                            padding: const EdgeInsets.symmetric(
                                                horizontal: 16),
                                            child: Icon(
                                              buttonMenu,
                                              color: CustomColors.whiteColorLow,
                                              size: 20,
                                            ),
                                          ),
                                        ),
                                        onTap: () {
                                          setState(
                                            () {
                                              extend = !extend;
                                              buttonMenu = extend
                                                  ? Icons.arrow_back
                                                  : Icons.menu;
                                            },
                                          );
                                        },
                                      ),
                                      Container(
                                        height: 40,
                                        width: 40,
                                        decoration: const BoxDecoration(
                                          image: DecorationImage(
                                            image: AssetImage(
                                                "assets/logo_white.png"),
                                          ),
                                        ),
                                      ),
                                    ],
                                  ),
                                  CustomProfile(
                                    name: profile[0],
                                    role: profile[1],
                                    message: toolTips[0],
                                    letter: letter,
                                    imageName: language.currentLanguage,
                                    onTap: () async {
                                      Completer<int> action = Completer<int>();
                                      showMenu(
                                        context: context,
                                        position: const RelativeRect.fromLTRB(
                                            60, 68, 0, 0),
                                        items: [
                                          PopupMenuItem(
                                            child: Row(
                                              mainAxisAlignment:
                                                  MainAxisAlignment
                                                      .spaceBetween,
                                              children: [
                                                Text(
                                                  language.getDataLanguage(
                                                          language
                                                              .currentLanguage)[
                                                      "dictionary"]["edit"],
                                                  style: const TextStyle(
                                                      fontSize: 13),
                                                ),
                                                const Padding(
                                                  padding: EdgeInsets.only(
                                                      left: 8.0),
                                                  child: Icon(
                                                    Icons.edit_note,
                                                    color: CustomColors
                                                        .primaryColorApp,
                                                    size: 20,
                                                  ),
                                                )
                                              ],
                                            ),
                                            onTap: () {
                                              action.complete(0);
                                            },
                                          ),
                                          PopupMenuItem(
                                            child: Row(
                                              mainAxisAlignment:
                                                  MainAxisAlignment
                                                      .spaceBetween,
                                              children: [
                                                Text(
                                                  language.getDataLanguage(
                                                          language
                                                              .currentLanguage)[
                                                      "dictionary"]["exit"],
                                                  style: const TextStyle(
                                                      fontSize: 13),
                                                ),
                                                const Padding(
                                                  padding: EdgeInsets.only(
                                                      left: 8.0),
                                                  child: Icon(
                                                    Icons.exit_to_app_outlined,
                                                    color: Colors.redAccent,
                                                    size: 20,
                                                  ),
                                                )
                                              ],
                                            ),
                                            onTap: () {
                                              action.complete(1);
                                            },
                                          )
                                        ],
                                      );
                                      actionMenu(
                                          await action.future,
                                          profile,
                                          translator,
                                          homeProvider,
                                          language.currentLanguage,
                                          authService);
                                    },
                                  ),
                                ],
                              ),
                            ),
                          ),
                          Container(
                            height: constraints.maxHeight - heightHeader,
                            width: constraints.maxWidth,
                            color: CustomColors.whiteColorLow,
                            child: Row(
                              children: [
                                Container(
                                  width:
                                      extend ? widthButtonMax : widthButtonMin,
                                  height: constraints.maxHeight,
                                  color: CustomColors.primaryColorApp,
                                  child: ListView.builder(
                                    scrollDirection: Axis.vertical,
                                    itemCount: buttonIcons.length,
                                    itemBuilder: (context, index) =>
                                        // Container(
                                        //   height: 48,
                                        //   color: getRandomColor(),
                                        // )
                                        CustomButtonVerticalHome(
                                      buttonNames: toolTips.sublist(2),
                                      toolTips: toolTips[index + 2],
                                      index: index,
                                      icons: buttonIcons,
                                      currentPosition: currentPosition,
                                      heightButton: heightButton,
                                      minMenu: extend,
                                      onTap: () {
                                        setState(
                                          () {
                                            currentPosition = index;
                                            pageviewNavigator(
                                                pageController, index);
                                          },
                                        );
                                      },
                                    ),
                                  ),
                                ),
                                Expanded(
                                  child: Padding(
                                    padding: const EdgeInsets.all(8),
                                    child: MultiProvider(
                                      providers: [
                                        ChangeNotifierProvider(
                                            create: (_) => DeviceViewModel()),
                                        ChangeNotifierProvider(
                                            create: (_) => LogViewModel()),
                                        ChangeNotifierProvider(
                                            create: (_) => SettignsViewModel())
                                      ],
                                      child: PageView(
                                          physics:
                                              const NeverScrollableScrollPhysics(),
                                          controller: pageController,
                                          scrollDirection: Axis.vertical,
                                          children: const [
                                            PageViewDevices(),
                                            Log(),
                                            Settigns(),
                                          ]),
                                    ),
                                  ),
                                )
                              ],
                            ),
                          )
                        ],
                      ),
                    ),
                  );
                }
              },
            );
          },
        ),
      ),
    );
  }
}

void pageviewNavigator(PageController pageController, int position) {
  pageController.animateToPage(position,
      duration: const Duration(milliseconds: 1000), curve: Curves.easeIn);
}

Future<bool> editProfile(BuildContext context, List<String> profile,
    Map<String, dynamic> translator, AuthService authService) async {
  Completer<bool> action = Completer<bool>();
  bool isSwitched = false;
  AlertMessage alertMessage = AlertMessage();

  List<TextEditingController> controllerList = List.generate(
    4,
    (index) => TextEditingController(
        text: index == 0
            ? profile[0]
            : index == 1
                ? profile[1]
                : ""),
  );

  final List<String> tags = [
    translator["dictionary"]["user"],
    translator["dictionary"]["email"],
    translator["dictionary"]["password"],
    translator["dictionary"]["confirmPassword"],
  ];

  Widget textField(
      double width,
      int index,
      bool obscureText,
      String? error,
      Function(String)? onChanged,
    ) {
    return CustomProfileTextField(
      width: width,
      controller: controllerList[index],
      labelText: tags[index],
      enabled: index > 1 ? isSwitched : true,
      isRequired: true,
      isEmail: (index == 1),
      isUser: (index == 0),
      obscureText: obscureText,
      translator: translator,
      error: error,
      onChanged: onChanged,
    );
    // return SizedBox(
    //   width: width,
    //   child: Column(
    //     children: [
    //       TextField(
    //         enabled: index > 1 ? isSwitched : true,
    //         controller: controllerList[index],
    //         style: const TextStyle(fontSize: 13),
    //         inputFormatters: [
    //           if (index > 0) FilteringTextInputFormatter.deny(RegExp(r'\s')),
    //         ],
    //         decoration: InputDecoration(
    //           filled: true,
    //           fillColor: CustomColors.whiteColorHigh,
    //           labelText: tags[index],
    //           labelStyle: const TextStyle(color: Colors.black54, fontSize: 13),
    //           border: OutlineInputBorder(
    //             borderRadius: BorderRadius.circular(2),
    //             borderSide: BorderSide.none,
    //           ),
    //         ),
    //       ),
    //       item
    //     ],
    //   ),
    // );
  }

  // Widget textFieldAlert(double width, bool show, String message) {
  //   print("width: $width, show: $show, message: $message");
  //   return Container(
  //     height: 12,
  //     width: width,
  //     padding: const EdgeInsets.only(left: 8),
  //     alignment: Alignment.topLeft,
  //     child: Text(
  //       show ? translator["devices"][message] : "",
  //       style: const TextStyle(fontSize: 10, color: Colors.red),
  //     ),
  //   );
  // }

  void close() {
    Navigator.of(context).pop(action);
  }

  showDialog(
    context: context,
    barrierDismissible: false,
    builder: (BuildContext context) {
      bool iconPasswordState = true;
      bool iconPasswordHover = false;
      String error = '';
      return StatefulBuilder(
        builder: (context, setState) {
          return Consumer<ProfileProvider>(
            builder: (context, login, child) => Dialog(
              child: Padding(
                padding: const EdgeInsets.symmetric(horizontal: 16.0),
                child: SizedBox(
                  height: 420,
                  width: 580,
                  child: Column(
                    children: [
                      Container(
                        alignment: Alignment.centerLeft,
                        height: 80,
                        child: Text(
                          translator["dictionary"]["edit"],
                          style: const TextStyle(fontSize: 24),
                        ),
                      ),
                      Container(
                        height: 224,
                        width: 580,
                        padding: const EdgeInsets.only(top: 24),
                        child: Row(
                          crossAxisAlignment: CrossAxisAlignment.start,
                          children: [
                            Container(
                              height: 100,
                              width: 100,
                              alignment: Alignment.center,
                              decoration: BoxDecoration(
                                  color: CustomColors.secondaryColorApp,
                                  borderRadius: BorderRadius.circular(90)),
                              child: Text(
                                profile[0][0],
                                style: const TextStyle(
                                    fontSize: 48,
                                    color: CustomColors.whiteColorHigh),
                              ),
                            ),
                            Container(
                              padding: const EdgeInsets.only(left: 24),
                              height: 200,
                              width: 480,
                              child: Column(
                                children: [
                                  Wrap(
                                    spacing: 11.0,
                                    runSpacing: 16.0,
                                    alignment: WrapAlignment.start,
                                    children: List<Widget>.generate(
                                      tags.length + 2,
                                      (index) => index < 2
                                          ? textField(
                                              220,
                                              index,
                                              false,
                                              null,
                                              null,
                                            )
                                          : index == 2
                                              ? SizedBox(
                                                  height: 40,
                                                  width: 480,
                                                  child: Row(
                                                    children: [
                                                      Switch(
                                                          value: isSwitched,
                                                          onChanged: (value) {
                                                            setState(
                                                              () {
                                                                isSwitched =
                                                                    value;
                                                                if (!isSwitched) {
                                                                  controllerList[
                                                                          2]
                                                                      .text = "";
                                                                  controllerList[
                                                                          3]
                                                                      .text = "";
                                                                }
                                                              },
                                                            );
                                                          }),
                                                      Text(translator["devices"]
                                                          ["updatePassword"])
                                                    ],
                                                  ),
                                                )
                                              : (index == 3) ?
                                              MouseRegion(
                                                child: GestureDetector(
                                                  child: Container(
                                                    height: 48,
                                                    width: 48,
                                                    // margin: const EdgeInsets.all(4),
                                                    decoration: BoxDecoration(
                                                        borderRadius: BorderRadius.circular(90),
                                                        color: iconPasswordHover ? CustomColors.hoverColorLogin : null),
                                                    child: Icon(
                                                      iconPasswordState
                                                          ? Icons.visibility_off_outlined
                                                          : Icons.remove_red_eye_outlined,
                                                      size: 20,
                                                      color: isSwitched ? CustomColors.greyColorHigh : CustomColors.neutral200,
                                                    ),
                                                  ),
                                                  onTap: () {
                                                    if (isSwitched) {
                                                      setState(() {
                                                        iconPasswordState =
                                                        !iconPasswordState;
                                                      });
                                                    }
                                                  },
                                                ),
                                                onEnter: (event) {
                                                  if (isSwitched) {
                                                    setState(() {
                                                      iconPasswordHover = true;
                                                    });
                                                  }
                                                },
                                                onExit: (event) {
                                                  if (isSwitched) {
                                                    setState(() {
                                                      iconPasswordHover = false;
                                                    });
                                                  }
                                                },
                                              )
                                              :
                                              textField(
                                                190,
                                                index - 2,
                                                iconPasswordState,
                                                index == 5 ? error : null,
                                                (text) {
                                                  error = '';
                                                  setState(() {
                                                    error;
                                                  });
                                                },
                                              ),
                                    ),
                                  ),
                                ],
                              ),
                            )
                          ],
                        ),
                      ),
                      SizedBox(
                        height: 110,
                        width: 580,
                        child: Row(
                          mainAxisAlignment: MainAxisAlignment.end,
                          children: [
                            Container(
                              width: 100,
                              alignment: Alignment.center,
                              child: TextButton(
                                onPressed: () {
                                  action.complete(false);
                                  close();
                                },
                                child: Text(
                                  translator["devices"]["cancel"],
                                  style: const TextStyle(
                                      fontSize: 13, color: Colors.black),
                                ),
                              ),
                            ),
                            Consumer<LanguageProvider>(
                              builder: (context, language, child) =>
                                  CustomRoundedButton(
                                textName: translator["devices"]["save"],
                                height: 32,
                                width: 100,
                                fontSize: 12,
                                isSelected: true,
                                textColorActived: CustomColors.whiteColorLow,
                                textColorInactive: CustomColors.whiteColorLow,
                                splashColor: CustomColors.primaryColorApp,
                                backgroundColorActived:
                                    CustomColors.primaryColorApp,
                                backgroundColorInactive:
                                    CustomColors.primaryColorApp,
                                borderRadiusValue: 30,
                                onTap: () async {
                                  if (controllerList[0].text == "") {
                                    login.updateNameState(true);
                                  } else {
                                    login.updateNameState(false);
                                  }

                                  if (controllerList[1].text == "") {
                                    login.updateEmailState(true);
                                  } else {
                                    login.updateEmailState(false);
                                  }

                                  if (controllerList[2].text !=
                                      controllerList[3].text) {
                                    login.updatePasswordState(true);
                                  } else {
                                    login.updatePasswordState(false);
                                  }

                                  if (controllerList[2].text.isNotEmpty &&
                                      controllerList[3].text.isNotEmpty) {
                                    login.updatePasswordState(false);
                                  } else {
                                    login.updatePasswordState(true);
                                  }
                                  bool save = false;
                                  late String password = "";
                                  if (isSwitched && !login.currentPassword) {
                                    if (controllerList[2].text ==
                                        controllerList[3].text) {
                                      password = controllerList[3].text;
                                      save = true;
                                    } else {
                                      login.updatePasswordState(true);
                                      save = false;
                                    }
                                  } else {
                                    if (!isSwitched &&
                                        !login.currentName &&
                                        !login.currentEmail) {
                                      save = true;
                                    }
                                  }

                                  save = false;
                                  bool invalidEmailMatch = !RegExp(r'^[\w-\.]+@([\w-]+\.)+[\w-]{2,4}$').hasMatch(controllerList[1].text);

                                  if (controllerList[0].text != profile[0] ||
                                      controllerList[1].text != profile[1]) {
                                    if (controllerList[0].text.length > 2 &&
                                        !invalidEmailMatch) {
                                      save = true;
                                    } else {
                                      if (controllerList[0].text.length < 3) {
                                        if (context.mounted) {
                                          alertMessage.showError(
                                              context,
                                              translator["dictionary"]['msg_low_user']);
                                        }
                                      }
                                      if (invalidEmailMatch) {
                                        if (context.mounted) {
                                          alertMessage.showError(
                                              context,
                                              translator["dictionary"]['msg_email_wrong']);
                                        }
                                      }
                                    }
                                  } else {
                                    if (!isSwitched) {
                                      close();
                                    }
                                  }

                                  if (isSwitched) {
                                    if (controllerList[2].text != "" &&
                                        controllerList[2].text == controllerList[3].text) {
                                      save = true;
                                    } else {
                                      if (controllerList[2].text == "") {
                                        if (context.mounted) {
                                          alertMessage.showError(
                                              context,
                                              translator["dictionary"]['msg_password_empty']);
                                        }
                                      }
                                      if (controllerList[3].text == "") {
                                        if (context.mounted) {
                                          alertMessage.showError(
                                              context,
                                              translator["dictionary"]['msg_confirm_password_empty']);
                                        }
                                      }
                                      if (controllerList[2].text != controllerList[3].text) {
                                        error = translator["dictionary"]['msg_passwords_not_match'];
                                        setState(() {
                                          error;
                                        });
                                        if (context.mounted) {
                                          alertMessage.showError(
                                              context,
                                              translator["dictionary"]['msg_passwords_not_match']);
                                        }
                                      }
                                    }
                                  }

                                  if (save) {
                                    int code = await authService.updateProfile(
                                        profile[2],
                                        controllerList[0].text,
                                        controllerList[1].text,
                                        password,
                                        language.currentLanguage);
                                    if (code == 200) {
                                      action.complete(true);
                                      close();
                                      if (context.mounted) {
                                        alertMessage.showSuccess(
                                            context,
                                            language.getDataLanguage(language
                                                .currentLanguage)["api"]
                                            ['edit_profile_200']);
                                      }
                                    } else {
                                      if (context.mounted) {
                                        alertMessage.showError(
                                            context,
                                            language.getDataLanguage(language
                                                .currentLanguage)["api"]
                                            [code.toString()]);
                                      }
                                    }
                                  }
                                },
                              ),
                            ),
                          ],
                        ),
                      ),
                    ],
                  ),
                ),
              ),
            ),
          );
        },
      );
    },
  );

  return action.future;
}
