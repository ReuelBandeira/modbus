import 'dart:convert';

import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:http/http.dart';
import 'package:isi_4_0/colors/color.dart';
import 'package:isi_4_0/providers/cards.dart';
import 'package:isi_4_0/providers/lang.dart';
import 'package:isi_4_0/repository/db_api.dart';
import 'package:isi_4_0/utils/alert_message.dart';
import 'package:isi_4_0/utils/data_devices.dart';
import 'package:isi_4_0/utils/data_protocols.dart';
import 'package:isi_4_0/utils/make_config.dart';
import 'package:isi_4_0/utils/mask.dart';
import 'package:isi_4_0/viewmodel/device_view_model.dart';
import 'package:isi_4_0/widgets/custom_dropbutton_vertical_with_button.dart';
import 'package:isi_4_0/widgets/custom_rounded_button.dart';
import 'package:provider/provider.dart';

class SaveConfiguration extends StatefulWidget {
  const SaveConfiguration({super.key});

  @override
  State<SaveConfiguration> createState() => _SaveConfigurationState();
}

class _SaveConfigurationState extends State<SaveConfiguration> {
  AuthService authService = AuthService.instance;
  AlertMessage alertMessage = AlertMessage();

  late List<DataProtocols> protocolList = [];
  late List<DataProtocols> newProtocolList = [];

  List<bool> alerts = List<bool>.generate(5, (index) => false);

  late DataDevices addedDevice;
  late List<String> topicList = [];
  late List<String> selectedItems = [];

  String? selectedProtocol;
  String? selectedTopic;

  List<TextEditingController> controllerList =
      List.generate(4, (index) => TextEditingController());

  TextEditingController topics = TextEditingController();
  TextEditingController topicsDialog = TextEditingController();

  final double heightTextField = 62;
  final double widthTextField = 180;

  void selectProtocol(String value) {
    setState(() {
      selectedProtocol = value;
    });
  }

  void selectTopic(String? value) {
    setState(() {
      selectedTopic = value;
    });
  }

  @override
  void initState() {
    super.initState();
    updatePage();
  }

  updatePage() async {
    var getProtocols = await authService.getProtocols();

    for (var protocol in getProtocols) {
      int id = protocol["id"];
      int version = protocol["version"];
      String name = protocol["protocol"];
      String alias = protocol["alias"];
      Map<String, dynamic> config = protocol["config"] == ""
          ? jsonDecode("{}")
          : jsonDecode(protocol["config"]);
      DataProtocols dataProtocol = DataProtocols(
          protocol: name, alias: alias, version: version, id: id, config: config);

      newProtocolList.add(dataProtocol);
    }

    setState(() {
      protocolList = newProtocolList;
    });
  }

  @override
  Widget build(BuildContext context) {
    void clearFields() {
      setState(() {
        selectedProtocol = null;
        selectedTopic = null;
        protocolList = [];
        topics.clear();
        topicList.clear();
        for (var controller in controllerList) {
          controller.clear();
        }
      });
    }

    Map<String, dynamic> translator = (context).select(
        (LanguageProvider lang) => lang.getDataLanguage(lang.currentLanguage));

    final List<String> deviceFields = [
      translator["devices"]["ipAddress"],
      translator["devices"]["port"],
      translator["devices"]["selectProtocol"],
      translator["devices"]["name"],
      translator["devices"]["update"],
      translator["devices"]["listTopics"],
    ];

    Widget textField(int index, Widget item) {
      return Column(
        children: [
          TextField(
            controller: controllerList[index],
            style: const TextStyle(fontSize: 13),
            inputFormatters: [
              if (index == 0) FilteringTextInputFormatter.deny(RegExp(r'\s')),
              if (index == 1 || index == 3)
                FilteringTextInputFormatter.digitsOnly
            ],
            decoration: InputDecoration(
              filled: true,
              fillColor: CustomColors.whiteColorHigh,
              labelText:
                  index < 2 ? deviceFields[index] : deviceFields[index + 1],
              labelStyle: const TextStyle(color: Colors.black54, fontSize: 13),
              border: OutlineInputBorder(
                borderRadius: BorderRadius.circular(2),
                borderSide: BorderSide.none,
              ),
            ),
          ),
          item
        ],
      );
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
        margin: const EdgeInsets.only(top: 4.0),
        height: MediaQuery.of(context).size.width > 700 ? 324 : 403,
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
                  translator["devices"]["title"],
                  style: const TextStyle(
                      fontSize: 16, fontWeight: FontWeight.bold),
                ),
                const SizedBox(height: 16.0),
                Text(
                  translator["devices"]["message"],
                  style: const TextStyle(
                      fontSize: 12, fontWeight: FontWeight.normal),
                ),
                const SizedBox(height: 16.0),
                Wrap(
                  spacing: 11.0,
                  runSpacing: 16.0,
                  alignment: WrapAlignment.start,
                  children: List<Widget>.generate(
                    deviceFields.length,
                    (index) => Container(
                      height: heightTextField,
                      width: widthTextField,
                      alignment: Alignment.center,
                      child: index == 2
                          ? SizedBox(
                              height: heightTextField,
                              width: widthTextField,
                              child: Column(
                                children: [
                                  Container(
                                    height: heightTextField - 12,
                                    width: widthTextField,
                                    color: CustomColors.whiteColorHigh,
                                    padding: const EdgeInsets.only(
                                        left: 16.0, right: 10.0),
                                    child: DropdownButton(
                                      items: protocolList.map(
                                        (DataProtocols option) {
                                          return DropdownMenuItem(
                                            value: option.protocol,
                                            child: Text(
                                              option.protocol,
                                              style:
                                                  const TextStyle(fontSize: 13),
                                            ),
                                          );
                                        },
                                      ).toList(),
                                      onChanged: (value) {
                                        setState(() {
                                          selectedProtocol = value.toString();
                                        });
                                      },
                                      underline: const SizedBox.shrink(),
                                      value: selectedProtocol,
                                      hint: Text(
                                        translator["devices"]["selectProtocol"],
                                        style: const TextStyle(
                                            color: Colors.black54,
                                            fontSize: 13),
                                      ),
                                      isExpanded: true,
                                    ),
                                  ),
                                  textFieldAlert(widthTextField, alerts[index],
                                      "msg_alert_field")
                                ],
                              ),
                            )
                          : index == 5
                              ? SizedBox(
                                  height: heightTextField,
                                  width: widthTextField,
                                  child: Column(
                                    children: [
                                      DropbuttonVerticalWithButton(
                                        heightContainer: heightTextField - 12,
                                        widthContainer: widthTextField,
                                        widthDropdownButton:
                                            widthTextField - 40,
                                        widthButton: 40,
                                        backgroundColor: CustomColors.whiteColorHigh,
                                        dropList: topicList,
                                        onChanged: (values) {},
                                        hintDrop: deviceFields[index],
                                        paddingLeft: 16.00,
                                        tooltipsMessage: translator["tooltips"]
                                            ["topicsMenu"],
                                        onPressedIcon: () {
                                          showCustomDialog(context);
                                        },
                                      ),
                                      textFieldAlert(widthTextField, false,
                                          "msg_alert_field")
                                    ],
                                  ),
                                )
                              : index < 5
                                  ? index < 2
                                      ? textField(
                                          index,
                                          textFieldAlert(widthTextField,
                                              alerts[index], "msg_alert_field"))
                                      : textField(
                                          index - 1,
                                          textFieldAlert(widthTextField,
                                              alerts[index], "msg_alert_field"))
                                  : Container(),
                    ),
                  ),
                ),
                Container(
                  height: 95,
                  width: constraints.maxWidth > 590 ? 604 : 2 * widthTextField,
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
                      Consumer<CardProvider>(
                        builder: (context, card, child) => CustomRoundedButton(
                          textName: translator["devices"]["save"],
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
                          onTap: () async {
                            JsonData jsonData = JsonData(
                                address: controllerList[0].text,
                                port: controllerList[1].text,
                                name: controllerList[2].text,
                                protocolConnection: selectedProtocol,
                                readingTime: controllerList[3].text,
                                topics: topicList,
                                data: [{}]);
                            // selectedProtocol == protocolList[2].protocol
                            //     ? modbus
                            //     : selectedProtocol ==
                            //             protocolList[1].protocol
                            //         ? ethernetIP
                            //         : canopen);

                            Map<String, dynamic> jsonMap = jsonData.toJson();

                            jsonMap.forEach((key, value) {
                              if (key == "address") {
                                alerts[0] =
                                    MaskDetector.detectMaskType(value) == "IP"
                                        ? false
                                        : true;
                              }
                              if (key == "port") {
                                alerts[1] = value == "" ? true : false;
                              }
                              if (key == "protocol") {
                                alerts[2] = value == null ? true : false;
                              }
                              if (key == "name") {
                                alerts[3] = value == "" ? true : false;
                              }
                              if (key == "readingTime") {
                                alerts[4] = value == "" ? true : false;
                              }
                            });

                            // Update Card
                            setState(() {});

                            // Check fields
                            int sum = alerts.fold(0,
                                (int previousValue, bool element) {
                              if (element == true) {
                                return previousValue + 1;
                              } else {
                                return previousValue;
                              }
                            });

                            // Save configuration
                            if (sum == 0) {
                              Response response = (await authService
                                  .saveDevice(jsonEncode(jsonMap))) as Response;
                              int code = response.statusCode;

                              if (code == 200) {
                                if (context.mounted) {
                                  alertMessage.showSuccess(
                                      context,
                                      translator["api"][code.toString()]);
                                }
                                clearFields();
                                card.updateConfiguredDevices();
                                var deviceProvider = DeviceViewModel();
                                dynamic jsonDataAdded =
                                    json.decode(response.body);
                                String id = jsonDataAdded["id"].toString();
                                List<String> topics =
                                    List<String>.from(jsonDataAdded["topics"]);
                                List<String> fields = [
                                  jsonDataAdded['address'],
                                  jsonDataAdded['port'],
                                  jsonDataAdded['protocol'],
                                  jsonDataAdded['name'],
                                  jsonDataAdded['readingTime'].toString(),
                                ];
                                List<dynamic> extraFields = [];
                                List<dynamic> extraFieldsData = jsonDataAdded['data'];

                                for (var item in newProtocolList) {
                                  if (item.protocol ==
                                      jsonDataAdded['protocol']) {
                                    for (var object in item.config.entries) {
                                      fields.add('');
                                      extraFields.add(object);
                                    }
                                  }
                                }

                                DataDevices dataDevices = DataDevices(
                                    fields: fields,
                                    topics: topics,
                                    id: id,
                                    extraFields: extraFields,
                                    extraFieldsData: extraFieldsData);

                                // Update data
                                deviceProvider.updateDataDevices(dataDevices);
                                // Call page
                                if (context.mounted) {
                                  deviceProvider.devicesExtraData(context);
                                }
                              } else {
                                if (context.mounted) {
                                  alertMessage.showError(
                                      context,
                                      translator["api"][code.toString()]);
                                }
                                clearFields();
                              }
                            }
                          },
                        ),
                      )
                    ],
                  ),
                ),
              ],
            ),
          ),
        ),
      );
    });
  }

  void showCustomDialog(BuildContext context) {
    showDialog(
      context: context,
      barrierDismissible: false,
      builder: (BuildContext context) {
        return Consumer<LanguageProvider>(
          builder: (context, language, child) => StatefulBuilder(
            builder: (context, setState) {
              return Dialog(
                child: Padding(
                  padding: const EdgeInsets.symmetric(horizontal: 16.0),
                  child: SizedBox(
                    height: 480,
                    width: 600,
                    child: Column(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: [
                        Container(
                          alignment: Alignment.centerLeft,
                          height: 80,
                          child: Text(
                            language.getDataLanguage(
                                    language.currentLanguage)["devices"]
                                ["configTopics"],
                            style: const TextStyle(fontSize: 24),
                          ),
                        ),
                        SizedBox(
                          height: 80,
                          child: Row(
                            children: [
                              Container(
                                width: 280,
                                padding: const EdgeInsets.only(right: 8.0),
                                child: TextField(
                                  controller: topicsDialog,
                                  style: const TextStyle(fontSize: 13),
                                  inputFormatters: [
                                    FilteringTextInputFormatter.deny(
                                      RegExp(r'\s'),
                                    ),
                                  ],
                                  decoration: InputDecoration(
                                    filled: true,
                                    fillColor: CustomColors.whiteColorHigh,
                                    labelText: language.getDataLanguage(
                                            language.currentLanguage)["devices"]
                                        ["topic"],
                                    labelStyle: const TextStyle(
                                        color: Colors.black54, fontSize: 13),
                                    border: OutlineInputBorder(
                                      borderRadius: BorderRadius.circular(2),
                                      borderSide: BorderSide.none,
                                    ),
                                  ),
                                ),
                              ),
                              Tooltip(
                                message: language.getDataLanguage(
                                        language.currentLanguage)["devices"]
                                    ["topicAddition"],
                                decoration: BoxDecoration(
                                  color: CustomColors.primaryColorApp.withOpacity(0.75),
                                  borderRadius: BorderRadius.circular(3),
                                ),
                                child: SizedBox(
                                  height: 50,
                                  child: TextButton(
                                    child:
                                        const Icon(Icons.playlist_add_outlined),
                                    onPressed: () {
                                      setState(() {
                                        topicList.add(topicsDialog.text);
                                        topicsDialog.clear();
                                      });
                                    },
                                  ),
                                ),
                              )
                            ],
                          ),
                        ),
                        Container(
                          height: 220,
                          width: 600,
                          color: CustomColors.whiteColorHigh,
                          child: SingleChildScrollView(
                            child: Padding(
                              padding: const EdgeInsets.all(8.0),
                              child: Wrap(
                                alignment: WrapAlignment.start,
                                spacing: 8,
                                runSpacing: 16,
                                children: List<Widget>.generate(
                                  topicList.length,
                                  (topic) => GestureDetector(
                                    onTap: () {
                                      showDialog(
                                        context: context,
                                        builder: (BuildContext context) {
                                          return AlertDialog(
                                            title: Text(language
                                                    .getDataLanguage(language
                                                        .currentLanguage)[
                                                "devices"]["removeTopic"]),
                                            content: Text(
                                                '${language.getDataLanguage(language.currentLanguage)["devices"]["alertTopic"]}"${topicList[topic]}" ?'),
                                            actions: [
                                              Padding(
                                                padding:
                                                    const EdgeInsets.all(8.0),
                                                child: Row(
                                                  mainAxisAlignment:
                                                      MainAxisAlignment.end,
                                                  children: [
                                                    Padding(
                                                      padding:
                                                          const EdgeInsets.only(
                                                              right: 16.0),
                                                      child: TextButton(
                                                        onPressed: () {
                                                          Navigator.of(context)
                                                              .pop();
                                                        },
                                                        child: Text(
                                                          language.getDataLanguage(
                                                                  language
                                                                      .currentLanguage)[
                                                              "devices"]["cancel"],
                                                          style:
                                                              const TextStyle(
                                                                  fontSize: 13,
                                                                  color: Colors
                                                                      .black),
                                                        ),
                                                      ),
                                                    ),
                                                    CustomRoundedButton(
                                                      textName: language
                                                              .getDataLanguage(
                                                                  language
                                                                      .currentLanguage)[
                                                          "devices"]["remove"],
                                                      backgroundColorActived:
                                                          CustomColors.primaryColorApp,
                                                      backgroundColorInactive:
                                                          CustomColors.primaryColorApp,
                                                      isSelected: true,
                                                      height: 32,
                                                      width: 100,
                                                      fontSize: 12,
                                                      splashColor:
                                                          CustomColors.primaryColorApp,
                                                      textColorActived:
                                                          CustomColors.whiteColorLow,
                                                      textColorInactive:
                                                          CustomColors.whiteColorLow,
                                                      borderRadiusValue: 30,
                                                      onTap: () async {
                                                        setState(
                                                          () {
                                                            topicList.removeAt(
                                                                topic);
                                                            Navigator.of(
                                                                    context)
                                                                .pop();
                                                          },
                                                        );
                                                      },
                                                    ),
                                                  ],
                                                ),
                                              ),
                                            ],
                                          );
                                        },
                                      );
                                    },
                                    onDoubleTap: () {
                                      TextEditingController
                                          textEditingController =
                                          TextEditingController(
                                              text: topicList[topic]);
                                      showDialog(
                                        context: context,
                                        builder: (BuildContext context) {
                                          return AlertDialog(
                                            title: Text(language
                                                    .getDataLanguage(language
                                                        .currentLanguage)[
                                                "devices"]["editTopic"]),
                                            content: TextField(
                                              controller: textEditingController,
                                              inputFormatters: [
                                                FilteringTextInputFormatter
                                                    .deny(RegExp(r'\s')),
                                              ],
                                            ),
                                            actions: [
                                              Padding(
                                                padding:
                                                    const EdgeInsets.all(8.0),
                                                child: Row(
                                                  mainAxisAlignment:
                                                      MainAxisAlignment
                                                          .spaceBetween,
                                                  children: [
                                                    Padding(
                                                      padding:
                                                          const EdgeInsets.only(
                                                              left: 16.0),
                                                      child: TextButton(
                                                        onPressed: () {
                                                          Navigator.of(context)
                                                              .pop();
                                                        },
                                                        child: Text(
                                                          language.getDataLanguage(
                                                                  language
                                                                      .currentLanguage)[
                                                              "devices"]["cancel"],
                                                          style:
                                                              const TextStyle(
                                                                  fontSize: 13,
                                                                  color: Colors
                                                                      .black),
                                                        ),
                                                      ),
                                                    ),
                                                    CustomRoundedButton(
                                                      textName: language
                                                              .getDataLanguage(
                                                                  language
                                                                      .currentLanguage)[
                                                          "devices"]["save"],
                                                      backgroundColorActived:
                                                          CustomColors.primaryColorApp,
                                                      backgroundColorInactive:
                                                          CustomColors.primaryColorApp,
                                                      isSelected: true,
                                                      height: 32,
                                                      width: 100,
                                                      fontSize: 12,
                                                      splashColor:
                                                          CustomColors.primaryColorApp,
                                                      textColorActived:
                                                          CustomColors.whiteColorLow,
                                                      textColorInactive:
                                                          CustomColors.whiteColorLow,
                                                      borderRadiusValue: 30,
                                                      onTap: () async {
                                                        setState(
                                                          () {
                                                            topicList[topic] =
                                                                textEditingController
                                                                    .text;
                                                            Navigator.of(
                                                                    context)
                                                                .pop();
                                                          },
                                                        );
                                                      },
                                                    ),
                                                  ],
                                                ),
                                              ),
                                            ],
                                          );
                                        },
                                      );
                                    },
                                    child: InputChip(
                                      disabledColor:
                                          CustomColors.primaryColorApp.withOpacity(0.4),
                                      label: Text(
                                        topicList[topic],
                                        style: const TextStyle(
                                            color: Colors.black),
                                      ),
                                    ),
                                  ),
                                ),
                              ),
                            ),
                          ),
                        ),
                        SizedBox(
                          height: 100,
                          child: Row(
                            mainAxisAlignment: MainAxisAlignment.end,
                            children: [
                              CustomRoundedButton(
                                textName: language.getDataLanguage(
                                        language.currentLanguage)["devices"]
                                    ["close"],
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
                                  if (topicList.isNotEmpty) {
                                    selectTopic(topicList.first);
                                  } else {
                                    selectTopic(null);
                                  }
                                  Navigator.of(context).pop();
                                },
                              )
                            ],
                          ),
                        ),
                      ],
                    ),
                  ),
                ),
              );
            },
          ),
        );
      },
    );
  }
}
