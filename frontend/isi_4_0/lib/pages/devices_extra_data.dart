import 'dart:async';
import 'dart:convert';

import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:isi_4_0/colors/color.dart';
import 'package:isi_4_0/utils/data_devices.dart';
import 'package:isi_4_0/widgets/custom_rounded_button.dart';
import 'package:provider/provider.dart';

import '../providers/lang.dart';
import '../repository/db_api.dart';
import '../utils/alert_message.dart';
import '../utils/data_protocols.dart';
import '../utils/make_config.dart';
import '../utils/mask.dart';
import '../widgets/custom_dropbutton_vertical_with_button.dart';

class DevicesExtraData extends StatefulWidget {
  final DataDevices dataDevices;

  const DevicesExtraData({super.key, required this.dataDevices});

  @override
  State<DevicesExtraData> createState() => _DevicesExtraDataState();
}

class _DevicesExtraDataState extends State<DevicesExtraData> {
  AlertMessage alertMessage = AlertMessage();
  AuthService authService = AuthService.instance;

  late List<DataProtocols> protocolList = [];
  late List<String> topicList = [];
  late dynamic extraDataObj;
  late dynamic extraDataControllerMap = {};

  @override
  void initState() {
    super.initState();
    updatePage();
  }

  updatePage() async {
    extraDataObj = widget.dataDevices.extraFieldsData[0];
    setState(() {
      extraDataObj;
    });

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
  }

  @override
  Widget build(BuildContext context) {
    List<TextEditingController> controllerList = List.generate(
        widget.dataDevices.fields.length,
        (index) =>
            TextEditingController(text: widget.dataDevices.fields[index]));

    List<bool> alerts = List<bool>.generate(50, (index) => false);

    List<String> columnTitle = [
      "IP address",
      "Port",
      "Protocol",
      "Name",
      "Update Interval",
      "Topics"
    ];

    const double heightTextField = 162;
    const double widthTextField = 180;

    const double heightAlert = 12;

    String? selectedProtocol;

    List<String> topics = widget.dataDevices.topics;
    TextEditingController topicsDialog = TextEditingController();

    final List<String> listaImutavel =
        List<String>.unmodifiable(List<String>.from(topics));

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
              : widget.dataDevices.extraFields[index - 6].value['name'],
          labelStyle: const TextStyle(color: Colors.black54, fontSize: 13),
          border: OutlineInputBorder(
              borderRadius: BorderRadius.circular(2),
              borderSide: BorderSide.none),
        ),
      );
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

    Widget buildDynamicField(
        dynamic extraField,
        dynamic extraDataObj,
        dynamic extraDataObjValue,
        dynamic extraDataObjKey,
        dynamic controllerMap) {
      // developer.log(
      //   'buildDynamicField',
      //   name: 'buildDynamicField value',
      //   error: widget.dataDevices.extraFieldsData[0].values.toList()[index].toString(),
      // );
      // developer.log(
      //   'extraField item',
      //   name: 'extraField item value',
      //   error: extraField.toString(),
      // );
      // developer.log(
      //   'extraDataObjValue item',
      //   name: 'extraDataObjValue item value',
      //   error: extraDataObjValue.toString(),
      // );

      if ((extraField.value['type'] == "string") ||
          (extraField.value['type'] == "int") ||
          (extraField.value['type'] == "float")) {
        TextEditingController controller = TextEditingController(
            text:
                extraDataObjValue == null ? '' : extraDataObjValue.toString());
        controllerMap[extraDataObjKey.toString()] = {
          'controller': controller,
          'type': extraField.value['type']
        };
        // controllerList.add(controller);
        return TextField(
          controller: controller,
          style: const TextStyle(fontSize: 13),
          inputFormatters: [
            FilteringTextInputFormatter.deny(
              RegExp(r'\s'),
            ),
            if ((extraField.value['type'] == "int") ||
                (extraField.value['type'] == "float"))
              FilteringTextInputFormatter.digitsOnly,
          ],
          decoration: InputDecoration(
            filled: true,
            fillColor: CustomColors.whiteColorHigh,
            labelText: extraField.value['name'],
            labelStyle: const TextStyle(color: Colors.black54, fontSize: 13),
            border: OutlineInputBorder(
                borderRadius: BorderRadius.circular(2),
                borderSide: BorderSide.none),
          ),
        );
      } else {
        if (extraField.value['type'] == "array") {
          controllerMap[extraDataObjKey.toString()] = {
            'controller': [],
            'type': extraField.value['type']
          };
          dynamic mainObj = {};
          for (var object in extraField.value.entries) {
            if (object.key == 'contains') {
              mainObj = object;
            }
          }
          List<Widget> data() {
            List<Widget> list = List.empty(growable: true);

            if (extraDataObjValue != null) {
              for (var i = 0; i < extraDataObjValue.length; i++) {
                controllerMap[extraDataObjKey.toString()]['controller'].add({});
                Widget listItem = Column(
                  children: [
                    Card(
                      child: buildDynamicField(
                          mainObj,
                          extraDataObjValue,
                          extraDataObjValue[i],
                          extraDataObjKey,
                          controllerMap[extraDataObjKey.toString()]
                              ['controller'][i]),
                    ),
                    ElevatedButton(
                      style: ElevatedButton.styleFrom(
                        primary: Colors.lightBlue,
                        padding: const EdgeInsets.all(12),
                        textStyle: const TextStyle(fontSize: 22),
                      ),
                      onPressed: () {
                        extraDataObjValue.removeAt(i);
                        updatePage();
                      },
                      child: const Text('Remove'),
                    )
                  ],
                );
                list.add(listItem);
              }
            }

            return list;
          }

          return ListView(children: [
            Text(extraField.key.toString()),
            SingleChildScrollView(
                scrollDirection: Axis.vertical,
                child: Card(
                    child: Column(
                  children: [
                    ...data(),
                    ElevatedButton(
                      style: ElevatedButton.styleFrom(
                        primary: Colors.lightBlue,
                        padding: const EdgeInsets.all(12),
                        textStyle: const TextStyle(fontSize: 22),
                      ),
                      onPressed: () {
                        if (extraDataObjValue == null) {
                          extraDataObj[extraDataObjKey] = [];
                          extraDataObjValue = extraDataObj[extraDataObjKey];
                        }
                        extraDataObjValue.add({});
                        updatePage();
                      },
                      child: const Text('Add'),
                    )
                  ],
                )))
          ]);
        } else {
          // object
          controllerMap[extraDataObjKey.toString()] = {
            'controller': {},
            'type': extraField.value['type']
          };
          List<Widget> data() {
            List<Widget> list = List.empty(growable: true);
            for (var object in extraField.value.entries) {
              if (object.key == 'contains') {
                for (var objectField in object.value.entries) {
                  // developer.log(
                  //   'extraDataObjValue item',
                  //   name: 'extraDataObjValue item value',
                  //   error: extraDataObjValue.toString(),
                  // );
                  // developer.log(
                  //   'objectField.key item',
                  //   name: 'objectField.key item value',
                  //   error: objectField.key.toString(),
                  // );
                  list.add(Card(
                    child: buildDynamicField(
                        objectField,
                        extraDataObjValue,
                        extraDataObjValue[objectField.key],
                        objectField.key,
                        controllerMap[extraDataObjKey.toString()]
                            ['controller']),
                  ));
                }
              }
            }
            return list;
          }

          return Column(
            children: data(),
          );
        }
      }
    }

    // Reset states from  Main Dialog with Statefull Father
    void resetState() {
      topics.clear();
      setState(() {
        topics.addAll(listaImutavel);
      });
    }

    Map<String, dynamic> fetchDynamicControllerText(dynamic entries,
        [bool parentIsArray = false]) {
      Map<String, dynamic> extraData = {};

      for (var item in entries) {
        // developer.log(
        //   'extraDataControllerMap item',
        //   name: 'extraDataControllerMap item value',
        //   error: item.toString(),
        // );
        if (item.value['controller'] is TextEditingController) {
          // developer.log(
          //   'extraDataControllerMap item is controller',
          //   name: 'extraDataControllerMap item is controller value',
          //   error: item.value.text.toString(),
          // );
          if (item.value['type'] == "int") {
            extraData[item.key] = int.parse(item.value['controller'].text);
          } else if (item.value['type'] == "float") {
            extraData[item.key] = double.parse(item.value['controller'].text);
          } else {
            extraData[item.key] = item.value['controller'].text;
          }
        } else if (item.value['controller'] is List) {
          // developer.log(
          //   'extraDataControllerMap item is array',
          //   name: 'extraDataControllerMap item is array value',
          //   error: item.value.toString(),
          // );
          extraData[item.key] = [];
          for (var arrayItem in item.value['controller']) {
            // developer.log(
            //   'extraDataControllerMap item is array and in array',
            //   name: 'extraDataControllerMap item is array  and in array value',
            //   error: arrayItem.toString(),
            // );
            if (arrayItem['controller'] is TextEditingController) {
              if (item.value['type'] == "int") {
                extraData[item.key]
                    .add(int.parse(arrayItem['controller'].text));
              } else if (item.value['type'] == "float") {
                extraData[item.key]
                    .add(double.parse(arrayItem['controller'].text));
              } else {
                extraData[item.key].add(arrayItem['controller'].text);
              }
            } else {
              extraData[item.key]
                  .add(fetchDynamicControllerText(arrayItem.entries, true));
            }
          }
        } else {
          // object
          // developer.log(
          //   'extraDataControllerMap item is object',
          //   name: 'extraDataControllerMap item is object value',
          //   error: item.value.toString(),
          // );
          if (parentIsArray) {
            extraData =
                fetchDynamicControllerText(item.value['controller'].entries);
          } else {
            extraData[item.key] =
                fetchDynamicControllerText(item.value['controller'].entries);
          }
          // for (var objItem in item.value.entries) {
          //   developer.log(
          //     'extraDataControllerMap item is object and in object',
          //     name: 'extraDataControllerMap item is object  and in object value',
          //     error: objItem.toString(),
          //   );
          //   extraData[item.key][objItem.key] =
          //       fetchDynamicControllerText(item.value.entries);
          // }
        }
      }

      return extraData;
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

      for (int i = 0; i < controllerList.length; i++) {
        if (i == 2) {
          fields.add(selectedProtocol ?? device.fields[2]);
        } else {
          fields.add(controllerList[i].text);
        }
      }

      extraDataList
          .add(fetchDynamicControllerText(extraDataControllerMap.entries));

      // print(extraDataList.toString());

      // if (device.fields[2] == "modbus") {
      //   extraDataList[0]["bitMemories"] = [
      //     {"address": 40960, "name": "Y00"},
      //     {"address": 40961, "name": "Y01"},
      //     {"address": 40962, "name": "Y02"},
      //     {"address": 40963, "name": "Y03"},
      //   ];
      //   // extraDataList[0]["slaveId"] = 14;
      //   extraDataList[0]["wordMemories"] = [
      //     {"format": 16, "address": 0, "name": "D0"},
      //     {"format": 16, "address": 1, "name": "D1"},
      //     {"format": 16, "address": 2, "name": "D2"},
      //     {"format": 16, "address": 3, "name": "D3"},
      //   ];
      // }
      // else if (device.fields[2] == "EthernetIP") {
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
      // extraDataList.add(extraData);

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

    return MaterialApp(
      home: Scaffold(
        body: Container(
          height: double.infinity,
          width: double.infinity,
          color: getRandomColor(),
          child: StatefulBuilder(
            builder: (context, setState) {
              return Consumer<LanguageProvider>(
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
                            const Text(
                              "Edit device",
                              style: TextStyle(
                                  fontSize: 18, fontWeight: FontWeight.bold),
                            ),
                            Tooltip(
                              message: "Remove topic",
                              decoration: BoxDecoration(
                                color: CustomColors.primaryColorApp.withOpacity(0.75),
                                borderRadius: BorderRadius.circular(3),
                              ),
                              child: SizedBox(
                                height: heightTextField,
                                width: heightTextField,
                                child: TextButton(
                                  onPressed: () async {
                                    if (await removeAction(context,
                                        widget.dataDevices.fields[3], 1)) {
                                      if (await authService.deleteDevice(
                                          widget.dataDevices.id)) {
                                        updatePage();
                                        alertMessage.message(
                                            "Device removed successfully!",
                                            Colors.green,
                                            '#4caf50');
                                        Navigator.pop(context);
                                      } else {
                                        alertMessage.message(
                                            "Device removal attempt failed!",
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
                          widget.dataDevices.fields.length + 1,
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
                                                widget.dataDevices.fields[2],
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
                                                    "Topic editing",
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
                                                                          const Text(
                                                                        "Edit topics",
                                                                        style: TextStyle(
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
                                                                                labelText: "Topic",
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
                                                                                "Add topic",
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
                                                                                  if (await removeAction(context, topics[topic], 0)) {
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
                                                                                        title: const Text("Edit topics"),
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
                                                                                                      Navigator.pop(context);
                                                                                                    },
                                                                                                    child: const Text(
                                                                                                      "Cancel",
                                                                                                      style: TextStyle(fontSize: 13, color: Colors.black),
                                                                                                    ),
                                                                                                  ),
                                                                                                ),
                                                                                                CustomRoundedButton(
                                                                                                  textName: "SAVE",
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
                                                                                                        Navigator.pop(context);
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
                                                                                "CLOSE",
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
                                                                              Navigator.pop(context);
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
                                              child: buildDynamicField(
                                                  widget.dataDevices
                                                      .extraFields[index - 6],
                                                  extraDataObj,
                                                  extraDataObj[widget
                                                      .dataDevices
                                                      .extraFields[index - 6]
                                                      .key],
                                                  widget
                                                      .dataDevices
                                                      .extraFields[index - 6]
                                                      .key,
                                                  extraDataControllerMap),
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
                                Navigator.pop(context);
                              },
                              child: const Text(
                                "CANCEL",
                                style: TextStyle(
                                    fontSize: 13, color: Colors.black),
                              ),
                            ),
                            const SizedBox(width: 16),
                            Padding(
                              padding: const EdgeInsets.only(right: 4.0),
                              child: CustomRoundedButton(
                                textName: "SAVE",
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
                                      widget.dataDevices,
                                      topics,
                                      language)) {
                                  } else {
                                    Navigator.pop(context);
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
              );
            },
          ),
        ),
      ),
    );
  }
}

Future<bool> removeAction(BuildContext context, String name, int id) async {
  Completer<bool> action = Completer<bool>();

  late String title;
  late String message;
  late String cancel;
  late String confirm;

  cancel = "CANCEL";
  confirm = "Remove";

  switch (id) {
    case 0:
      title = "Remove topic";
      message = '${"Are you sure you want to remove the topic "}"$name" ?';
      break;
    case 1:
      title = "Remove device";
      message = '${"Are you sure you want to remove the device "}"$name" ?';
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
                      Navigator.pop(context);
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
                    Navigator.pop(context);
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
