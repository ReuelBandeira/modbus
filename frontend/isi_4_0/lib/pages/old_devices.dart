import 'dart:async';
import 'dart:convert';

import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:isi_4_0/colors/color.dart';
import 'package:isi_4_0/providers/cards.dart';
import 'package:isi_4_0/providers/lang.dart';
import 'package:isi_4_0/repository/db_api.dart';
import 'package:isi_4_0/utils/alert_message.dart';
import 'package:isi_4_0/utils/data_devices.dart';
import 'package:isi_4_0/utils/data_protocols.dart';
import 'package:isi_4_0/utils/language.dart';
import 'package:isi_4_0/utils/make_config.dart';
import 'package:isi_4_0/utils/mask.dart';
import 'package:isi_4_0/viewmodel/device_view_model.dart';
import 'package:isi_4_0/widgets/custom_card_horizontal_home.dart';
import 'package:isi_4_0/widgets/custom_dropbutton_vertical_with_button.dart';
import 'package:isi_4_0/widgets/custom_rounded_button.dart';
import 'package:provider/provider.dart';

class OldDevices extends StatefulWidget {
  const OldDevices({super.key});

  @override
  State<OldDevices> createState() => _OldDevicesState();
}

class _OldDevicesState extends State<OldDevices> {
  AuthService authService = AuthService.instance;
  AlertMessage alertMessage = AlertMessage();

  List<bool> alerts = List<bool>.generate(50, (index) => false);

  late List<String> dataCard = List.filled(4, "");

  late final Function(String) atualizarValor;

  late List<DataDevices> dataList = [];
  late List<DataProtocols> protocolList = [];
  late List<String> topicList = [];
  late List<String> topicListAux = [];

  final double widthLine = 1206;
  final double heightAlert = 12;

  @override
  void initState() {
    super.initState();
    updatePage();
  }

  @override
  void dispose() {
    super.dispose();
  }

  updatePage() async {
    try {
      var countDevice = await authService.countDevices();
      var getDevices = await authService.getDevices(1, 10);
      var getProtocols = await authService.getProtocols();

      protocolList.clear();

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

        protocolList.add(dataProtocol);
      }

      if (countDevice != "0") {
        dataList.clear();

        for (var deviceList in getDevices['items']) {
          var devices = deviceList["devices"];
          for (var device in devices) {
            String id = device["id"].toString();
            List<String> topics = List<String>.from(device["topics"]);
            List<String> fields = [
              device['address'],
              device['port'],
              device['protocol'],
              device['name'],
              device['readingTime'].toString(),
            ];
            List<dynamic> extraFields = [];
            List<dynamic> extraFieldsData = device['data'];

            for (var item in protocolList) {
              if (item.protocol == device['protocol']) {
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

            dataList.add(dataDevices);
          }
        }


      } else {
        setState(() {
          dataList = [];
        });
      }
      update();
    } catch (e) {
      print('Error on function: updatePage in Devices => $e');
    }
  }

  void update() {
    setState(() {});
  }

  openMainDialog(BuildContext context, int index, final DataDevices device,
      Map<String, dynamic> tagFields, List<String> columnTitle) async {
    alerts.map((e) => false);
    List<TextEditingController> controllerList = List.generate(
        device.fields.length,
        (index) => TextEditingController(text: device.fields[index]));

    TextEditingController topicsDialog = TextEditingController();

    List<String> topics = device.topics;

    final List<String> listaImutavel =
        List<String>.unmodifiable(List<String>.from(topics));

    String? selectedProtocol;

    const double heightTextField = 62;
    const double widthTextField = 180;

    // Reset states from  Main Dialog with Statefull Father
    void resetState() {
      topics.clear();
      setState(() {
        topics.addAll(listaImutavel);
      });
    }

    Widget textFieldAlert(double width, bool show, String message) {
      return Consumer<LanguageProvider>(
        builder: (context, language, child) => Container(
          height: heightAlert,
          width: width,
          padding: const EdgeInsets.only(left: 8),
          alignment: Alignment.topLeft,
          child: Text(
            show
                ? language.getDataLanguage(language.currentLanguage)["devices"]
                    [message]
                : "",
            style: const TextStyle(fontSize: 10, color: Colors.red),
          ),
        ),
      );
    }

    Future<bool> closeDialog(
        BuildContext context,
        List<TextEditingController> controllerList,
        String? selectedProtocol,
        DataDevices device,
        List<String> topics,
        LanguageProvider language) async {
      List<String> fields = [];
      Map<String, dynamic> translator =
          language.getDataLanguage(language.currentLanguage);

      List<Map<String, dynamic>> extraDataList = [];
      final Map<String, dynamic> extraData = {};

      for (int i = 0; i < controllerList.length; i++) {
        if (i == 2) {
          fields.add(selectedProtocol ?? device.fields[2]);
        } else if (i > 4) {
          String key = device.extraFields[i - 5].key;
          extraData[key] = controllerList[i].text;
        } else {
          fields.add(controllerList[i].text);
        }
      }

      // if (device.fields[2] == "Modbus") {
      //   extraData["bitMemories"] = [
      //     {"address": 40960, "name": "Y00"},
      //     {"address": 40961, "name": "Y01"},
      //     {"address": 40962, "name": "Y02"},
      //     {"address": 40963, "name": "Y03"},
      //   ];
      //   extraData["wordMemories"] = [
      //     {"format": 16, "address": 0, "name": "D0"},
      //     {"format": 16, "address": 1, "name": "D1"},
      //     {"format": 16, "address": 2, "name": "D2"},
      //     {"format": 16, "address": 3, "name": "D3"},
      //   ];
      // } else if (device.fields[2] == "EthernetIP") {
      //   extraData["bitMemories"] = [
      //     {"attribute": 0, "class": 849, "instance": 1, "name": "Y00"},
      //     {"attribute": 1, "class": 849, "instance": 1, "name": "Y01"},
      //     {"attribute": 2, "class": 849, "instance": 1, "name": "Y02"},
      //     {"attribute": 3, "class": 849, "instance": 1, "name": "Y03"},
      //   ];
      //   extraData["wordMemories"] = [
      //     {"attribute": 0, "class": 850, "instance": 2, "name": "D0"},
      //     {"attribute": 1, "class": 850, "instance": 2, "name": "D1"},
      //     {"attribute": 2, "class": 850, "instance": 2, "name": "D2"},
      //     {"attribute": 3, "class": 850, "instance": 2, "name": "D3"},
      //   ];
      // }
      extraDataList.add(extraData);

      JsonData jsonData = JsonData(
          address: fields[0],
          port: fields[1],
          name: fields[3],
          protocolConnection: fields[2],
          readingTime: fields[4],
          topics: topics,
          data: extraDataList);

      Map<String, dynamic> json = jsonData.toJson();

      json.forEach((key, value) {
        if (key == "plcAddress") {
          alerts[0] = MaskDetector.detectMaskType(value) == "IP" ? false : true;
        }
        if (key == "plcPort") {
          alerts[1] = value == "" ? true : false;
        }
        if (key == "protocol") {
          alerts[2] = value == null ? true : false;
        }
        if (key == "plcName") {
          alerts[3] = value == "" ? true : false;
        }
        if (key == "readingTime") {
          alerts[4] = value.toString() == '' ? true : false;
        }
      });

      // Check fields
      int sum = alerts.fold(0, (int previousValue, bool element) {
        if (element == true) {
          return previousValue + 1;
        } else {
          return previousValue;
        }
      });

      // Save configuration
      if (sum == 0) {
        int code = await authService.updateDevice(device.id, jsonEncode(json));

        if (code == 200) {
          alertMessage.serverMessage(code, translator["api"][code.toString()]);
          updatePage();
        } else {
          alertMessage.serverMessage(code, translator["api"][code.toString()]);
          updatePage();
        }
        return false;
      } else {
        return true;
      }
    }

    showDialog(
      context: context,
      barrierDismissible: false,
      builder: (BuildContext context) {
        Widget textField(int index,
            [bool digitsOnly = false, bool ipOnly = false]) {
          return TextField(
            controller:
                index < 6 ? controllerList[index] : controllerList[index - 1],
            style: const TextStyle(fontSize: 13),
            inputFormatters: [
              FilteringTextInputFormatter.deny(
                RegExp(r'\s'),
              ),
              if (digitsOnly)
                FilteringTextInputFormatter.digitsOnly
              else if (ipOnly)
                FilteringTextInputFormatter.allow(RegExp(r'^[\d.]+')),
            ],
            decoration: InputDecoration(
              filled: true,
              fillColor: CustomColors.whiteColorHigh,
              labelText: index < 6
                  ? columnTitle[index]
                  : device.extraFields[index - 6].value['name'],
              labelStyle: const TextStyle(color: Colors.black54, fontSize: 13),
              border: OutlineInputBorder(
                  borderRadius: BorderRadius.circular(2),
                  borderSide: BorderSide.none),
            ),
          );
        }

        return StatefulBuilder(
          builder: (context, setState) {
            return Dialog(
              child: Consumer<LanguageProvider>(
                builder: (context, language, child) => Container(
                  padding: const EdgeInsets.all(16),
                  height: MediaQuery.of(context).size.width > 700 ? 415 : 480,
                  width: 600,
                  decoration: BoxDecoration(
                      color: Colors.white,
                      borderRadius: BorderRadius.circular(4)),
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.center,
                    children: [
                      Container(
                        height: heightTextField,
                        padding: const EdgeInsets.symmetric(horizontal: 8.0),
                        child: Row(
                          mainAxisAlignment: MainAxisAlignment.spaceBetween,
                          children: [
                            Text(
                              tagFields["devices"]["msgEdit"],
                              style: const TextStyle(
                                  fontSize: 18, fontWeight: FontWeight.bold),
                            ),
                            Tooltip(
                              message: tagFields["devices"]["removeTopic"],
                              decoration: BoxDecoration(
                                color: CustomColors.primaryColorApp.withOpacity(0.75),
                                borderRadius: BorderRadius.circular(3),
                              ),
                              child: SizedBox(
                                height: heightTextField,
                                width: heightTextField,
                                child: TextButton(
                                  onPressed: () async {
                                    if (await removeAction(context, tagFields,
                                        device.fields[3], 1)) {
                                      if (await authService
                                          .deleteDevice(device.id)) {
                                        updatePage();
                                        alertMessage.message(
                                            tagFields["devices"]
                                                ["msg_remove_device_ok"],
                                            Colors.green,
                                            '#4caf50');
                                        Navigator.of(context).pop();
                                      } else {
                                        alertMessage.message(
                                            tagFields["devices"]
                                                ["msg_remove_device_fail"],
                                            Colors.red,
                                            '#dc1c13');
                                      }
                                    }
                                  },
                                  child: const Icon(
                                    Icons.delete_forever,
                                    color: Colors.red,
                                  ),
                                ),
                              ),
                            ),
                          ],
                        ),
                      ),
                      const SizedBox(height: 16.0),
                      Wrap(
                        spacing: 8.0,
                        runSpacing: 12.0,
                        alignment: WrapAlignment.start,
                        children: List<Widget>.generate(
                          device.fields.length + 1,
                          (index) {
                            return Container(
                              height: heightTextField,
                              width: widthTextField,
                              alignment: Alignment.topCenter,
                              child: index == 2
                                  ? SizedBox(
                                      height: heightTextField - 2,
                                      width: widthTextField,
                                      child: Column(
                                        children: [
                                          Container(
                                            height: heightTextField - 14,
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
                                                      style: const TextStyle(
                                                          fontSize: 13),
                                                    ),
                                                  );
                                                },
                                              ).toList(),
                                              onChanged: (value) {
                                                setState(() {
                                                  selectedProtocol =
                                                      value.toString();
                                                });
                                              },
                                              underline:
                                                  const SizedBox.shrink(),
                                              value: selectedProtocol,
                                              hint: Text(
                                                device.fields[2],
                                                style: const TextStyle(
                                                    color: Colors.black,
                                                    fontSize: 13),
                                              ),
                                              isExpanded: true,
                                            ),
                                          ),
                                          textFieldAlert(widthTextField,
                                              alerts[index], "msg_alert_field")
                                        ],
                                      ),
                                    )
                                  : index == 5
                                      ? SizedBox(
                                          height: heightTextField - 2,
                                          width: widthTextField,
                                          child: Column(
                                            children: [
                                              DropbuttonVerticalWithButton(
                                                heightContainer:
                                                    heightTextField - 14,
                                                widthContainer: widthTextField,
                                                widthDropdownButton:
                                                    widthTextField - 40,
                                                widthButton: 40,
                                                backgroundColor:
                                                    CustomColors.whiteColorHigh,
                                                dropList: topics,
                                                onChanged: (values) {},
                                                hintDrop: columnTitle[5],
                                                tooltipsMessage:
                                                    tagFields["tooltips"]
                                                        ["topicEditing"],
                                                paddingLeft: 12.00,
                                                onPressedIcon: () async {
                                                  showDialog(
                                                    context: context,
                                                    barrierDismissible: false,
                                                    builder:
                                                        (BuildContext context) {
                                                      return StatefulBuilder(
                                                        builder: (context,
                                                            setState) {
                                                          return Dialog(
                                                            child: Padding(
                                                              padding: const EdgeInsets
                                                                      .symmetric(
                                                                  horizontal:
                                                                      16.0),
                                                              child: SizedBox(
                                                                height: 480,
                                                                width: 600,
                                                                child: Column(
                                                                  crossAxisAlignment:
                                                                      CrossAxisAlignment
                                                                          .start,
                                                                  children: [
                                                                    Container(
                                                                      alignment:
                                                                          Alignment
                                                                              .centerLeft,
                                                                      height:
                                                                          80,
                                                                      child:
                                                                          Text(
                                                                        tagFields["devices"]
                                                                            [
                                                                            "editTopic"],
                                                                        style: const TextStyle(
                                                                            fontSize:
                                                                                24),
                                                                      ),
                                                                    ),
                                                                    SizedBox(
                                                                      height:
                                                                          80,
                                                                      child:
                                                                          Row(
                                                                        children: [
                                                                          Container(
                                                                            width:
                                                                                280,
                                                                            padding:
                                                                                const EdgeInsets.only(right: 8.0),
                                                                            child:
                                                                                TextField(
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
                                                                                labelText: tagFields["devices"]["topic"],
                                                                                labelStyle: const TextStyle(color: Colors.black54, fontSize: 13),
                                                                                border: OutlineInputBorder(
                                                                                  borderRadius: BorderRadius.circular(2),
                                                                                  borderSide: BorderSide.none,
                                                                                ),
                                                                              ),
                                                                            ),
                                                                          ),
                                                                          Tooltip(
                                                                            message:
                                                                                tagFields["devices"]["topicAddition"],
                                                                            decoration:
                                                                                BoxDecoration(
                                                                              color: CustomColors.primaryColorApp.withOpacity(0.75),
                                                                              borderRadius: BorderRadius.circular(3),
                                                                            ),
                                                                            child:
                                                                                SizedBox(
                                                                              height: 50,
                                                                              child: TextButton(
                                                                                child: const Icon(Icons.playlist_add_outlined),
                                                                                onPressed: () {
                                                                                  setState(() {
                                                                                    topics.add(topicsDialog.text);
                                                                                    topicsDialog.clear();
                                                                                  });
                                                                                },
                                                                              ),
                                                                            ),
                                                                          ),
                                                                        ],
                                                                      ),
                                                                    ),
                                                                    Container(
                                                                      height:
                                                                          220,
                                                                      width:
                                                                          600,
                                                                      color:
                                                                          CustomColors.whiteColorHigh,
                                                                      child:
                                                                          SingleChildScrollView(
                                                                        child:
                                                                            Padding(
                                                                          padding:
                                                                              const EdgeInsets.all(8.0),
                                                                          child:
                                                                              Wrap(
                                                                            alignment:
                                                                                WrapAlignment.start,
                                                                            spacing:
                                                                                8,
                                                                            runSpacing:
                                                                                16,
                                                                            children:
                                                                                List.generate(
                                                                              topics.length,
                                                                              (topic) => GestureDetector(
                                                                                onTap: () async {
                                                                                  if (await removeAction(context, tagFields, topics[topic], 0)) {
                                                                                    setState(() {
                                                                                      topics.removeAt(topic);
                                                                                    });
                                                                                  }
                                                                                },
                                                                                onDoubleTap: () {
                                                                                  TextEditingController textEditingController = TextEditingController(text: topics[topic]);
                                                                                  showDialog(
                                                                                    context: context,
                                                                                    builder: (BuildContext context) {
                                                                                      return AlertDialog(
                                                                                        title: Text(tagFields["devices"]["editTopic"]),
                                                                                        content: TextField(
                                                                                          controller: textEditingController,
                                                                                          inputFormatters: [
                                                                                            FilteringTextInputFormatter.deny(RegExp(r'\s')),
                                                                                          ],
                                                                                        ),
                                                                                        actions: [
                                                                                          Padding(
                                                                                            padding: const EdgeInsets.all(8.0),
                                                                                            child: Row(
                                                                                              mainAxisAlignment: MainAxisAlignment.spaceBetween,
                                                                                              children: [
                                                                                                Padding(
                                                                                                  padding: const EdgeInsets.only(left: 16.0),
                                                                                                  child: TextButton(
                                                                                                    onPressed: () {
                                                                                                      Navigator.of(context).pop();
                                                                                                    },
                                                                                                    child: Text(
                                                                                                      tagFields["devices"]["cancel"],
                                                                                                      style: const TextStyle(fontSize: 13, color: Colors.black),
                                                                                                    ),
                                                                                                  ),
                                                                                                ),
                                                                                                CustomRoundedButton(
                                                                                                  textName: tagFields["devices"]["save"],
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
                                                                                                    setState(
                                                                                                      () {
                                                                                                        topics[topic] = textEditingController.text;
                                                                                                        Navigator.of(context).pop();
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
                                                                                  disabledColor: CustomColors.primaryColorApp.withOpacity(0.4),
                                                                                  label: Text(
                                                                                    topics[topic],
                                                                                    style: const TextStyle(color: Colors.black),
                                                                                  ),
                                                                                ),
                                                                              ),
                                                                            ),
                                                                          ),
                                                                        ),
                                                                      ),
                                                                    ),
                                                                    Container(
                                                                      height:
                                                                          70,
                                                                      alignment:
                                                                          Alignment
                                                                              .bottomCenter,
                                                                      child:
                                                                          Row(
                                                                        mainAxisAlignment:
                                                                            MainAxisAlignment.end,
                                                                        children: [
                                                                          CustomRoundedButton(
                                                                            textName:
                                                                                tagFields["devices"]["close"],
                                                                            backgroundColorActived:
                                                                                CustomColors.primaryColorApp,
                                                                            backgroundColorInactive:
                                                                                CustomColors.primaryColorApp,
                                                                            isSelected:
                                                                                true,
                                                                            height:
                                                                                32,
                                                                            width:
                                                                                100,
                                                                            fontSize:
                                                                                12,
                                                                            splashColor:
                                                                                CustomColors.primaryColorApp,
                                                                            textColorActived:
                                                                                CustomColors.whiteColorLow,
                                                                            textColorInactive:
                                                                                CustomColors.whiteColorLow,
                                                                            borderRadiusValue:
                                                                                30,
                                                                            onTap:
                                                                                () async {
                                                                              topicList = topicList;
                                                                              Navigator.of(context).pop(true);
                                                                            },
                                                                          ),
                                                                        ],
                                                                      ),
                                                                    ),
                                                                  ],
                                                                ),
                                                              ),
                                                            ),
                                                          );
                                                        },
                                                      );
                                                    },
                                                  ).then(
                                                    (action) {
                                                      if (action != null) {
                                                        if (action) {
                                                          setState(
                                                            () {},
                                                          );
                                                        }
                                                      }
                                                    },
                                                  );
                                                },
                                              ),
                                              textFieldAlert(widthTextField,
                                                  false, "msg_alert_field")
                                            ],
                                          ),
                                        )
                                      : index < 5
                                          ? SizedBox(
                                              height: heightTextField,
                                              width: widthTextField,
                                              child: Column(
                                                children: [
                                                  textField(
                                                      index,
                                                      ((index == 1) ||
                                                          (index == 4)),
                                                      (index == 0)),
                                                  textFieldAlert(
                                                      widthTextField,
                                                      alerts[index],
                                                      "msg_alert_field")
                                                ],
                                              ),
                                            )
                                          : SizedBox(
                                              height: heightTextField,
                                              width: widthTextField,
                                              child: Column(
                                                children: [
                                                  (device.extraFields[index - 6]
                                                                      .value[
                                                                  'type'] ==
                                                              "string") ||
                                                          (device
                                                                      .extraFields[
                                                                          index - 6]
                                                                      .value[
                                                                  'type'] ==
                                                              "int") ||
                                                          (device
                                                                      .extraFields[
                                                                          index - 6]
                                                                      .value[
                                                                  'type'] ==
                                                              "float")
                                                      ? textField(
                                                          index,
                                                          (device.extraFields[index - 6].value[
                                                                      'type'] ==
                                                                  "int") ||
                                                              (device
                                                                      .extraFields[
                                                                          index -
                                                                              6]
                                                                      .value['type'] ==
                                                                  "float"))
                                                      : Container(),
                                                  textFieldAlert(
                                                      widthTextField,
                                                      alerts[index],
                                                      "msg_alert_field")
                                                ],
                                              ),
                                            ),
                            );
                          },
                        ),
                      ),
                      Container(
                        margin: const EdgeInsets.only(top: 24.0),
                        height: 70,
                        alignment: Alignment.centerRight,
                        child: Row(
                          mainAxisAlignment: MainAxisAlignment.end,
                          children: [
                            TextButton(
                              onPressed: () async {
                                resetState();
                                Navigator.of(context).pop();
                              },
                              child: Text(
                                tagFields["devices"]["cancel"],
                                style: const TextStyle(
                                    fontSize: 13, color: Colors.black),
                              ),
                            ),
                            const SizedBox(width: 16),
                            Padding(
                              padding: const EdgeInsets.only(right: 4.0),
                              child: CustomRoundedButton(
                                textName: tagFields["devices"]["save"],
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
                                  if (await closeDialog(
                                      context,
                                      controllerList,
                                      selectedProtocol,
                                      device,
                                      topics,
                                      language)) {
                                  } else {
                                    Navigator.of(context).pop(true);
                                  }
                                },
                              ),
                            )
                          ],
                        ),
                      )
                    ],
                  ),
                ),
              ),
            );
          },
        );
      },
    );
  }

  @override
  Widget build(BuildContext context) {
    // Update cards
    (context).select((CardProvider card) => card.updateAll());
    Map<String, dynamic> translator = (context).select(
        (LanguageProvider lang) => lang.getDataLanguage(lang.currentLanguage));

    List<String> columnTitle = [
      translator["devices"]["ipAddress"],
      translator["devices"]["port"],
      translator["devices"]["title"],
      translator["devices"]["name"],
      translator["devices"]["update"],
      translator["devices"]["topics"]
    ];

    const double heightContainer = 48;

    return NotificationListener<OverscrollIndicatorNotification>(
      onNotification: (OverscrollIndicatorNotification overscroll) {
        overscroll.disallowIndicator();
        return true;
      },
      child: SizeChangedLayoutNotifier(
        child: Consumer<LanguageProvider>(
          builder: (context, language, child) {
            return LayoutBuilder(
              builder: (context, constraints) {
                // Warning: breaking language changes to old screens
                // List<String> cardTitle = Language()
                //     .getListFromCardTitleJson(language.currentLanguage);
                List<String> cardTitle = [];
                dataCard[0] = (context)
                    .select((CardProvider card) => card.configuredDevices);
                dataCard[1] =
                    (context).select((CardProvider card) => card.usedMemory);
                dataCard[2] =
                    (context).select((CardProvider card) => card.diskUsed);
                dataCard[3] = (context)
                    .select((CardProvider card) => card.connectedDevices);
                return Scaffold(
                  backgroundColor: CustomColors.whiteColorLow,
                  body: OverflowBox(
                    minHeight: 0,
                    minWidth: 0,
                    maxHeight: constraints.maxHeight,
                    maxWidth: constraints.maxWidth,
                    child: SingleChildScrollView(
                      child: Column(
                        crossAxisAlignment: CrossAxisAlignment.start,
                        children: [
                          SizedBox(
                            height: 140,
                            width: constraints.maxWidth,
                            child: ListView.builder(
                              scrollDirection: Axis.horizontal,
                              itemCount: cardTitle.length,
                              itemBuilder: (context, index) =>
                                  Consumer<CardProvider>(
                                builder: (context, card, child) =>
                                    CustomCardHorizontalHome(
                                  heightCard: 140,
                                  widthCard: 300,
                                  backgroundColor: CustomColors.whiteColorHigh,
                                  cardName: cardTitle[index],
                                  cardData: dataCard[index],
                                  cardIcon: Icons.cached_outlined,
                                  onTap: () {
                                    switch (index) {
                                      case 0:
                                        card.updateConfiguredDevices();
                                        break;
                                      case 1:
                                        card.updateUsedMemory();
                                        break;
                                      case 2:
                                        card.updateDiskUsed();
                                        break;
                                      case 3:
                                        card.updateConnectedDevices();
                                        break;
                                      default:
                                    }
                                  },
                                  shadowColor: CustomColors.greenColorLight,
                                ),
                              ),
                            ),
                          ),
                          // ========= ========= ========= Headers
                          Container(
                            height: heightContainer,
                            width: widthLine,
                            margin: const EdgeInsets.only(left: 6.0, top: 16),
                            child: ListView.builder(
                              scrollDirection: Axis.horizontal,
                              itemCount: columnTitle.length,
                              itemBuilder: (context, field) => Container(
                                height: heightContainer,
                                width: widthLine / columnTitle.length,
                                padding: const EdgeInsets.only(left: 16.0),
                                decoration: const BoxDecoration(
                                  color: Color.fromARGB(255, 240, 240, 240),
                                  border: Border(
                                    bottom: BorderSide(
                                      width: 4,
                                      color: Color.fromARGB(52, 57, 49, 49),
                                    ),
                                  ),
                                ),
                                alignment: Alignment.centerLeft,
                                child: Text(
                                  columnTitle[field],
                                  style: const TextStyle(
                                      fontWeight: FontWeight.bold),
                                ),
                              ),
                            ),
                          ),
                          // ========= ========= ========= Items - Colum
                          Container(
                            height: 10 * heightContainer,
                            width: widthLine,
                            margin: const EdgeInsets.only(left: 6.0),
                            child: ListView.builder(
                              scrollDirection: Axis.vertical,
                              itemCount: dataList.length,
                              itemBuilder: (context, device) {
                                return Container(
                                  height: heightContainer,
                                  width: widthLine,
                                  color: Colors.white,
                                  // ========= ========= ========= Items - Lines
                                  child: ListView.builder(
                                    scrollDirection: Axis.horizontal,
                                    // itemCount: dataList[device].fields.length + 1,
                                    itemCount: 6,
                                    itemBuilder: (context, field) {
                                      return Container(
                                        height: heightContainer,
                                        width: widthLine / 6,
                                        alignment: Alignment.centerLeft,
                                        padding:
                                            const EdgeInsets.only(left: 16.0),
                                        decoration: const BoxDecoration(
                                          color: Colors.white,
                                          border: Border(
                                            bottom: BorderSide(
                                              width: 1,
                                              color: CustomColors.whiteColorHigh,
                                            ),
                                          ),
                                        ),
                                        child: field < 5
                                            ? Text(
                                                dataList[device].fields[field])
                                            : Consumer<DeviceViewModel>(
                                                builder: (context,
                                                        deviceProvider,
                                                        child) =>
                                                    DropbuttonVerticalWithButton(
                                                  heightContainer:
                                                      heightContainer,
                                                  widthContainer:
                                                      (widthLine / 6),
                                                  widthDropdownButton:
                                                      (widthLine / 6) - 60,
                                                  widthButton: 44,
                                                  backgroundColor: Colors.white,
                                                  dropList:
                                                      dataList[device].topics,
                                                  onChanged: (values) {},
                                                  hintDrop:
                                                      translator["devices"]
                                                          ["listTopics"],
                                                  tooltipsMessage:
                                                      translator["tooltips"]
                                                          ["topicsMenu"],
                                                  onPressedIcon: () {
                                                    // Update data
                                                    deviceProvider
                                                        .updateDataDevices(
                                                            dataList[device]);
                                                    // Call page
                                                    deviceProvider
                                                        .devicesExtraData(
                                                            context);
                                                  },
                                                ),
                                              ),
                                      );
                                    },
                                  ),
                                );
                              },
                            ),
                          ),
                        ],
                      ),
                    ),
                  ),
                );
              },
            );
          },
        ),
      ),
    );
  }
}

Future<bool> removeAction(BuildContext context, Map<String, dynamic> tags,
    String name, int id) async {
  Completer<bool> action = Completer<bool>();

  late String title;
  late String message;
  late String cancel;
  late String confirm;

  cancel = tags["devices"]["cancel"];
  confirm = tags["devices"]["remove"];

  switch (id) {
    case 0:
      title = tags["devices"]["removeTopic"];
      message = '${tags["devices"]["alertTopic"]}"$name" ?';
      break;
    case 1:
      title = tags["devices"]["removeDevice"];
      message = '${tags["devices"]["alertDevice"]}"$name" ?';
      break;
    default:
  }

  showDialog(
    context: context,
    builder: (BuildContext context) {
      return AlertDialog(
        title: Text(title),
        content: Text(message),
        actions: [
          Padding(
            padding: const EdgeInsets.all(8.0),
            child: Row(
              mainAxisAlignment: MainAxisAlignment.end,
              children: [
                Padding(
                  padding: const EdgeInsets.only(right: 16.0),
                  child: TextButton(
                    onPressed: () {
                      action.complete(false);
                      Navigator.of(context).pop();
                    },
                    child: Text(
                      cancel,
                      style: const TextStyle(fontSize: 13, color: Colors.black),
                    ),
                  ),
                ),
                CustomRoundedButton(
                  textName: confirm,
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
                  onTap: () {
                    action.complete(true);
                    Navigator.of(context).pop();
                  },
                ),
              ],
            ),
          ),
        ],
      );
    },
  );

  return action.future;
}
