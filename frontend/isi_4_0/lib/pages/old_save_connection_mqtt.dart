import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:isi_4_0/providers/lang.dart';
import 'package:isi_4_0/utils/alert_message.dart';
import 'package:isi_4_0/colors/color.dart';
import 'package:isi_4_0/repository/db_api.dart';
import 'package:isi_4_0/utils/mask.dart';
import 'package:isi_4_0/widgets/custom_rounded_button.dart';
import 'package:provider/provider.dart';

class SaveConnectionMqtt extends StatefulWidget {
  const SaveConnectionMqtt({super.key});

  @override
  State<SaveConnectionMqtt> createState() => _SaveConnectionMqttState();
}

class _SaveConnectionMqttState extends State<SaveConnectionMqtt> {
  AuthService authService = AuthService.instance;

  List<TextEditingController> controllerList =
      List.generate(4, (index) => TextEditingController());

  List<bool> alerts = List<bool>.generate(4, (index) => false);

  AlertMessage alertMessage = AlertMessage();
  final double heightTextField = 62;
  final double widthTextField = 276;

  @override
  void initState() {
    super.initState();
  }

  @override
  Widget build(BuildContext context) {
    Map<String, dynamic> translator = (context).select(
        (LanguageProvider lang) => lang.getDataLanguage(lang.currentLanguage));

    final List<String> protocolFields = [
      translator["devices"]["ipAddress"],
      translator["devices"]["port"],
      translator["devices"]["name"],
      translator["dictionary"]["password"],
    ];

    void clearFields() {
      for (var controller in controllerList) {
        controller.clear();
      }
    }

    Widget textFieldAlert(double width, bool show, String message) {
      return Container(
        height: 12,
        width: width,
        padding: const EdgeInsets.only(left: 8),
        alignment: Alignment.topLeft,
        child: Text(
          show ? translator["devices"][message] : "",
          style: const TextStyle(fontSize: 10, color: Colors.red),
        ),
      );
    }

    return LayoutBuilder(
        builder: (BuildContext context, BoxConstraints constraints) {
      return Container(
        margin: const EdgeInsets.only(left: 4.0, top: 4.0),
        height: constraints.maxWidth > 590 ? 324 : 480,
        width: 604,
        child: Card(
          color: CustomColors.whiteColorLow,
          elevation: 2,
          shape: RoundedRectangleBorder(
            borderRadius: BorderRadius.circular(2.0),
          ),
          child: Padding(
            padding: const EdgeInsets.only(left: 16.0, top: 16.0),
            child: Column(
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
                  alignment: WrapAlignment.start,
                  children: List<Widget>.generate(
                    protocolFields.length,
                    (index) => Column(
                      children: [
                        Container(
                          height: heightTextField - 12,
                          width: widthTextField,
                          alignment: Alignment.center,
                          child: TextField(
                            controller: controllerList[index],
                            style: const TextStyle(fontSize: 13),
                            inputFormatters: [
                              FilteringTextInputFormatter.deny(
                                RegExp(r'\s'),
                              ),
                              if (index == 1)
                                FilteringTextInputFormatter.digitsOnly
                            ],
                            decoration: InputDecoration(
                              filled: true,
                              fillColor: CustomColors.whiteColorHigh,
                              labelText: protocolFields[index],
                              labelStyle: const TextStyle(
                                  color: Colors.black54, fontSize: 13),
                              border: OutlineInputBorder(
                                borderRadius: BorderRadius.circular(2),
                                borderSide: BorderSide.none,
                              ),
                            ),
                          ),
                        ),
                        textFieldAlert(
                            widthTextField, alerts[index], "msg_alert_field")
                      ],
                    ),
                  ),
                ),
                Container(
                  height: 95,
                  width: constraints.maxWidth > 590 ? 604 : widthTextField,
                  alignment: Alignment.center,
                  padding: const EdgeInsets.only(right: 16.0),
                  child: Row(
                    mainAxisAlignment: MainAxisAlignment.end,
                    children: [
                      Container(
                        width: 100,
                        alignment: Alignment.center,
                        child: TextButton(
                          onPressed: () {
                            clearFields();
                          },
                          child: Text(
                            translator["devices"]["cancel"],
                            style: const TextStyle(
                                fontSize: 13, color: Colors.black),
                          ),
                        ),
                      ),
                      CustomRoundedButton(
                        textName: translator["devices"]["save"],
                        backgroundColorActived: CustomColors.primaryColorApp,
                        backgroundColorInactive: CustomColors.primaryColorApp,
                        isSelected: true,
                        height: 32,
                        width: 100,
                        fontSize: 12,
                        splashColor: CustomColors.primaryColorApp,
                        textColorActived: CustomColors.whiteColorLow,
                        textColorInactive: CustomColors.whiteColorLow,
                        borderRadiusValue: 30,
                        onTap: () async {
                          alerts[0] = MaskDetector.detectMaskType(
                                      controllerList[0].text) ==
                                  "IP"
                              ? false
                              : true;
                          alerts[1] =
                              controllerList[1].text == "" ? true : false;

                          // Update Card
                          setState(() {});

                          // Check fields
                          int sum =
                              alerts.fold(0, (int previousValue, bool element) {
                            if (element == true) {
                              return previousValue + 1;
                            } else {
                              return previousValue;
                            }
                          });

                          if (sum == 0) {
                            int code = await authService.saveProtocolMQTT(
                                controllerList[0].text,
                                controllerList[1].text,
                                controllerList[2].text,
                                controllerList[3].text,
                                "inputinfo",
                                "errorinfo",
                                "statusdevice");
                            if (code == 200) {
                              alertMessage.serverMessage(
                                  code, translator["api"][code.toString()]);
                              clearFields();
                            } else {
                              alertMessage.serverMessage(
                                  code, translator["api"][code.toString()]);
                              clearFields();
                            }
                          }
                        },
                      ),
                    ],
                  ),
                )
              ],
            ),
          ),
        ),
      );
    });
  }
}
