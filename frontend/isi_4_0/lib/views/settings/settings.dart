import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:isi_4_0/colors/color.dart';
import 'package:isi_4_0/providers/lang.dart';
import 'package:isi_4_0/repository/db_api.dart';
import 'package:isi_4_0/utils/alert_message.dart';
import 'package:isi_4_0/utils/enum_login.dart';
import 'package:isi_4_0/utils/message_validator.dart';
import 'package:isi_4_0/utils/text_form_field.dart';
import 'package:isi_4_0/viewmodel/settings_view_model.dart';
import 'package:isi_4_0/widgets/custom_rounded_button.dart';
import 'package:provider/provider.dart';

class Settigns extends StatefulWidget {
  const Settigns({super.key});

  @override
  State<Settigns> createState() => _SettignsState();
}

class _SettignsState extends State<Settigns> {
  final formKeySettings = GlobalKey<FormState>();

  List<TextEditingController> controllerList =
      List.generate(4, (index) => TextEditingController());

  @override
  Widget build(BuildContext context) {
    return LayoutBuilder(builder: (context, mqttSize) {
      return SingleChildScrollView(
          scrollDirection: Axis.vertical,
          child: Consumer2<LanguageProvider, SettignsViewModel>(
              builder: (context, lang, settings, child) {
            return Column(
              mainAxisAlignment: MainAxisAlignment.start,
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                FutureBuilder<Map<String, dynamic>>(
                  future:
                      Future.value(lang.getDataLanguage(lang.currentLanguage)),
                  builder: (context, snapshot) {
                    if (snapshot.connectionState == ConnectionState.done) {
                      if (snapshot.hasData) {
                        return customCard(snapshot.data, controllerList,
                            formKeySettings, mqttSize, settings);
                      } else {
                        // Make Custom UX to Error Notifier
                        return Container();
                      }
                    } else {
                      // Make Custom loading...
                      return Container();
                    }
                  },
                )
              ],
            );
          }));
    });
  }
}

Widget customCard(
    Map<String, dynamic>? translator,
    List<TextEditingController> controllerList,
    GlobalKey<FormState> formKeySettings,
    BoxConstraints mqttSize,
    SettignsViewModel settings) {
  List<CurrentTextField> currentTextField = [
    CurrentTextField.addressIP,
    CurrentTextField.port,
    CurrentTextField.user,
    CurrentTextField.password
  ];

  final List<String> hints = [
    translator!["mqtt"]["ipAddress"],
    translator["mqtt"]["port"],
    translator["mqtt"]["name"],
    translator["mqtt"]["password"],
  ];
  const double heightTextField = 79;
  const double widthTextField = 276;

  return Form(
    key: formKeySettings,
    child: StatefulBuilder(builder: (context, update) {
      return Container(
        height: mqttSize.maxWidth > 635 ? 415 : 505,
        width: mqttSize.maxWidth > 635 ? 604 : 300,
        margin: const EdgeInsets.only(left: 32, top: 32),
        child: Card(
          color: CustomColors.whiteColorLow,
          elevation: 2,
          shape: RoundedRectangleBorder(
            borderRadius: BorderRadius.circular(4.0),
          ),
          child: Padding(
            padding: const EdgeInsets.all(16.0),
            child: SingleChildScrollView(
              scrollDirection: Axis.vertical,
              child: FutureBuilder<Map<String, dynamic>>(
                  future: settings.getMqttStatus(),
                  builder: (context, mqtt) {
                    if (mqtt.connectionState == ConnectionState.done) {
                      return Column(
                        crossAxisAlignment: CrossAxisAlignment.start,
                        children: [
                          Text(
                            translator["devices"]["mqtt"],
                            style: const TextStyle(
                                fontSize: 16, fontWeight: FontWeight.bold),
                          ),
                          const SizedBox(height: 16.0),
                          Text(
                            translator["devices"]["messageBroker"],
                            style: const TextStyle(
                                fontSize: 12, fontWeight: FontWeight.normal),
                          ),
                          const SizedBox(height: 16.0),
                          Wrap(
                            spacing: 11.0,
                            runSpacing: 16.0,
                            children: List.generate(
                              currentTextField.length,
                              (index) {
                                switch (index) {
                                  case 0:
                                    controllerList[index].text =
                                        mqtt.data!['server'];
                                    break;
                                  case 1:
                                    controllerList[index].text =
                                        mqtt.data!['port'];
                                    break;
                                  case 2:
                                    controllerList[index].text =
                                        mqtt.data!['username'];
                                    break;
                                  default:
                                }
                                return SizedBox(
                                  height: heightTextField,
                                  width: widthTextField,
                                  child: CustomTextFormField(
                                    controller: controllerList[index],
                                    currentTextField: currentTextField[index],
                                    currentPage: CurrentPage.settings,
                                    hintText: hints[index],
                                    inputFormatters: index == 1
                                        ? [
                                            FilteringTextInputFormatter
                                                .digitsOnly
                                          ]
                                        : null,
                                    onFieldSubmitted: (value) async {
                                      submitForm(context, formKeySettings,
                                          controllerList, translator);
                                    },
                                    validator: (value) {
                                      return MessageValidator.instance
                                          .getMessage(
                                              CurrentPage.settings,
                                              value,
                                              currentTextField[index],
                                              translator);
                                    },
                                  ),
                                );
                              },
                            ),
                          ),
                          Container(
                            height: 60,
                            padding: const EdgeInsets.only(left: 8.0),
                            alignment: Alignment.bottomLeft,
                            child: Row(
                              children: [
                                Text(
                                    '${translator["devices"]["mqttStatusLable"]}: '),
                                Text(
                                  mqtt.data!['connected']
                                      ? '${translator["devices"]["mqttOnline"]} '
                                      : '${translator["devices"]["mqttOffline"]} ',
                                  style: TextStyle(
                                      color: mqtt.data!['connected']
                                          ? CustomColors.success600
                                          : CustomColors.error600,
                                      fontWeight: FontWeight.bold),
                                ),
                                Padding(
                                  padding: const EdgeInsets.only(left: 8.0),
                                  child: CustomRoundedButton(
                                    textName:
                                        '${translator["devices"]["mqttTest"]}',
                                    height: 30,
                                    width: 120,
                                    fontSize: 12,
                                    isSelected: true,
                                    textColorActived:
                                        CustomColors.whiteColorLow,
                                    textColorInactive:
                                        CustomColors.whiteColorLow,
                                    splashColor: CustomColors.primaryColorApp,
                                    backgroundColorActived:
                                        CustomColors.primaryColorApp,
                                    backgroundColorInactive:
                                        CustomColors.primaryColorApp,
                                    borderRadiusValue: 30,
                                    onTap: () async {
                                      update(() {});
                                    },
                                  ),
                                )
                              ],
                            ),
                          ),
                          Container(
                            height: 75,
                            alignment: Alignment.bottomRight,
                            child: SizedBox(
                              height: 60,
                              width: 400,
                              child: Row(
                                mainAxisAlignment: MainAxisAlignment.end,
                                children: [
                                  SizedBox(
                                    height: 40,
                                    width: 100,
                                    child: TextButton(
                                      onPressed: () {
                                        clearFields(controllerList);
                                      },
                                      child: Text(
                                        translator["devices"]["cancel"],
                                        style: const TextStyle(
                                            fontSize: 13, color: Colors.black),
                                      ),
                                    ),
                                  ),
                                  Container(
                                    width: mqttSize.maxWidth > 635 ? 170 : 150,
                                    padding: const EdgeInsets.symmetric(
                                        horizontal: 16),
                                    child: CustomRoundedButton(
                                      textName: translator["devices"]["save"],
                                      height: 40,
                                      width: 100,
                                      fontSize: 12,
                                      isSelected: true,
                                      textColorActived:
                                          CustomColors.whiteColorLow,
                                      textColorInactive:
                                          CustomColors.whiteColorLow,
                                      splashColor: CustomColors.primaryColorApp,
                                      backgroundColorActived:
                                          CustomColors.primaryColorApp,
                                      backgroundColorInactive:
                                          CustomColors.primaryColorApp,
                                      borderRadiusValue: 30,
                                      onTap: () {
                                        submitForm(context, formKeySettings,
                                            controllerList, translator);
                                      },
                                    ),
                                  )
                                ],
                              ),
                            ),
                          )
                        ],
                      );
                    } else {
                      return Container();
                    }
                  }),
            ),
          ),
        ),
      );
    }),
  );
}

void clearFields(List<TextEditingController> controllerList) {
  for (var controller in controllerList) {
    controller.text = '';
  }
}

Future<void> submitForm(
    BuildContext context,
    GlobalKey<FormState> formKeySettings,
    List<TextEditingController> controllerList,
    Map<String, dynamic> translator) async {
  if (formKeySettings.currentState!.validate()) {
    AuthService authService = AuthService.instance;
    AlertMessage alertMessage = AlertMessage();
    int code = await authService.saveProtocolMQTT(
        controllerList[0].text,
        controllerList[1].text,
        controllerList[2].text,
        controllerList[3].text,
        "inputinfo",
        "errorinfo",
        "statusdevice");
    if (code == 200) {
      if (context.mounted) {
        alertMessage.showSuccess(context, translator["api"]['mqtt_200']);
      }
      clearFields(controllerList);
    } else {
      if (context.mounted) {
        alertMessage.showError(context, translator["api"][code.toString()]);
      }
    }
  }
}
