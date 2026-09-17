import 'dart:convert';

import 'package:flutter/material.dart';
import 'package:http/http.dart';
import 'package:isi_4_0/colors/color.dart';
import 'package:isi_4_0/providers/lang.dart';
import 'package:isi_4_0/views/devices/components/card_extra_data.dart';
import 'package:isi_4_0/views/devices/components/card_protocol.dart';
import 'package:isi_4_0/views/devices/components/card_textfield_protocol.dart';
import 'package:isi_4_0/views/devices/components/card_topics.dart';
import 'package:isi_4_0/widgets/custom_rounded_button.dart';
import 'package:provider/provider.dart';

import '../../repository/db_api.dart';
import '../../utils/alert_message.dart';
import '../../utils/data_protocols.dart';
import '../../utils/make_config.dart';
import '../../widgets/alert_blurry.dart';

class NewDevices extends StatefulWidget {
  const NewDevices(
      {super.key,
      this.onTap,
      this.extraData,
      required this.dataCallback,
      this.dataEdit});
  final Function()? onTap;
  final Function()? extraData;
  final void Function(dynamic) dataCallback;
  final dynamic dataEdit;
  @override
  State<NewDevices> createState() => _NewDevicesState();
}

class _NewDevicesState extends State<NewDevices> {
  AuthService authService = AuthService.instance;

  late List<DataProtocols> protocolList = [];
  late String? selectedProtocol;

  TextEditingController textControllerName = TextEditingController();
  TextEditingController textControllerIp = TextEditingController();
  TextEditingController textControllerPort = TextEditingController();
  TextEditingController textControllerReadTime = TextEditingController();

  late List<String> topicsList = [];

  late Widget? extraDataObj;
  late dynamic extraDataControllerMap = {};

  late List<Widget> extraDataRowList = [];

  late dynamic customData = {};

  AlertMessage alertMessage = AlertMessage();

  bool showExtraData = false;

  bool errorProtocol = false;
  bool errorConfiguration = false;
  bool errorTopics = false;
  bool errorExtraData = false;

  @override
  void initState() {
    extraDataObj = null;
    if (widget.dataEdit != null) {
      selectedProtocol = widget.dataEdit['protocol'];
      textControllerName.text = widget.dataEdit['name'];
      textControllerIp.text = widget.dataEdit['address'];
      textControllerPort.text = widget.dataEdit['port'];
      textControllerReadTime.text = widget.dataEdit['readingTime'].toString();

      topicsList = [];
      for (dynamic topic in widget.dataEdit['topics']) {
        topicsList.add(topic.toString());
      }

      customData = widget.dataEdit['data'][0];
    } else {
      selectedProtocol = null;
    }
    setState(() {
      selectedProtocol;
      protocolList;
      extraDataObj;
      topicsList;
      extraDataRowList;
      showExtraData;
      customData;
      extraDataControllerMap;
      errorProtocol;
      errorConfiguration;
      errorTopics;
      errorExtraData;
    });
    super.initState();
    updatePage();
    if (selectedProtocol != null) {
      loadExtraFields();
    }
  }

  @override
  void dispose() {
    super.dispose();
  }

  void rebuildAllChildren(BuildContext context) {
    void rebuild(Element el) {
      el.markNeedsBuild();
      el.visitChildren(rebuild);
    }
    (context as Element).visitChildren(rebuild);
  }

  void setShowExtraData(bool newState) {
    showExtraData = newState;
    setState(() {
      showExtraData;
    });
  }

  void setExtraDataRowList(List<Widget> newList) {
    extraDataRowList = newList;
    setState(() {
      extraDataRowList;
    });
  }

  List<Widget> getExtraDataRowList() {
    return extraDataRowList;
  }

  // void resetColumnsSelected(dynamic controllerMap) {
  //   if (controllerMap != null && controllerMap is! TextEditingController) {
  //     print(controllerMap);
  //     if ((controllerMap['select'] != null) && (controllerMap['hover'] != null)) {
  //       for (int i = 0; i < controllerMap['select'].length; i++) {
  //         controllerMap['select'][i] = false;
  //         controllerMap['hover'][i] = false;
  //       }
  //     }
  //
  //     if (controllerMap['controller'] != null) {
  //       if (controllerMap['controller'] is List) {
  //         for (dynamic controller in controllerMap['controller']) {
  //           for(dynamic entry in controller.entries) {
  //             resetColumnsSelected(entry.value);
  //           }
  //         }
  //       } else {
  //         resetColumnsSelected(controllerMap['controller']);
  //       }
  //     }
  //   }
  // }

  updatePage() async {
    // try {
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
          protocol: name,
          alias: alias,
          version: version,
          id: id,
          config: config);

      protocolList.add(dataProtocol);
    }

    setState(() {
      protocolList;
    });

    if (selectedProtocol != null) {
      loadExtraFields();
    }
    // } catch (e) {
    //   print('Error on function: updatePage in Devices => $e');
    // }
  }

  loadExtraFields() {
    CardExtraData card = CardExtraData(
      extraData: widget.extraData,
      extraDataObj: extraDataObj,
      setShowExtraData: setShowExtraData,
      setExtraDataRowList: setExtraDataRowList,
      getExtraDataRowList: getExtraDataRowList,
      isErrored: errorExtraData,
    );
    // try {
    for (DataProtocols protocol in protocolList) {
      if (protocol.protocol == selectedProtocol) {
        List<Widget> extraFields = [];
        for (var object in protocol.config.entries) {
          extraFields.add(card.buildDynamicField(
                  object,
                  customData,
                  customData[object.key],
                  object.key,
                  extraDataControllerMap,
                  widget.dataCallback,
                  true,
          ));
        }

        if (extraFields.isNotEmpty) {
          extraDataObj = Column(
            children: extraFields,
          );
          setState(() {
            extraDataObj;
          });
        }
        break;
      }
    }
    // } catch (e) {
    //   print('Error on function: loadExtraFields in Devices => $e');
    // }
  }

  Map<String, dynamic> fetchDynamicControllerText(
      dynamic mainControllerMap,
      dynamic entries,
      [bool parentIsArray = false]
    ) {
    Map<String, dynamic> extraData = {};

    for (var item in entries) {
      if (item.value['controller'] is TextEditingController) {
        if (item.value['required'] && item.value['controller'].text == '') {
          throw(mainControllerMap);
        }
        if (item.value['type'] == "int") {
          if (item.value['controller'].text == '') {
            extraData[item.key] = 0;
          } else {
            extraData[item.key] = int.parse(item.value['controller'].text);
          }
        } else if (item.value['type'] == "float") {
          if (item.value['controller'].text == '') {
            extraData[item.key] = 0;
          } else {
            extraData[item.key] = double.parse(item.value['controller'].text);
          }
        } else {
          extraData[item.key] = item.value['controller'].text;
        }
      } else if (item.value['controller'] is List) {
        extraData[item.key] = [];
        for (var arrayItem in item.value['controller']) {
          if (arrayItem['controller'] is TextEditingController) {
            if (item.value['required'] && item.value['controller'].text == '') {
              throw(mainControllerMap);
            }
            if (item.value['type'] == "int") {
              if (item.value['controller'].text == '') {
                extraData[item.key] = 0;
              } else {
                extraData[item.key].add(
                    int.parse(arrayItem['controller'].text));
              }
            } else if (item.value['type'] == "float") {
              if (item.value['controller'].text == '') {
                extraData[item.key] = 0;
              } else {
                extraData[item.key]
                    .add(double.parse(arrayItem['controller'].text));
              }
            } else {
              extraData[item.key].add(arrayItem['controller'].text);
            }
          } else {
            try {
              extraData[item.key]
                  .add(fetchDynamicControllerText(mainControllerMap, arrayItem.entries, true));
            } catch (e) {
              item.value['isErrored'] = true;
              throw(mainControllerMap);
            }
          }
        }
      } else {
        try {
          if (parentIsArray) {
            extraData =
                fetchDynamicControllerText(mainControllerMap, item.value['controller'].entries);
          } else {
            extraData[item.key] =
                fetchDynamicControllerText(mainControllerMap, item.value['controller'].entries);
          }
        } catch (e) {
          item.value['isErrored'] = true;
          throw(mainControllerMap);
        }
      }
    }

    return extraData;
  }

  Widget saveButton(Function()? onTap, BuildContext context) {
    Map<String, dynamic> translator = (context).select(
            (LanguageProvider lang) => lang.getDataLanguage(lang.currentLanguage));
    Future<bool> save(int? id) async {
      int code;
      errorProtocol = false;
      errorExtraData = false;
      errorTopics = false;
      errorConfiguration = false;
      setState(() {
        errorProtocol;
        errorExtraData;
        errorTopics;
        errorConfiguration;
      });
      if (selectedProtocol == null) {
        errorProtocol = true;
        setState(() {
          errorProtocol;
        });
        alertMessage.showError(context, translator["newDevices"]["errorSelectProtocol"]);
        return false;
      }
      if ((textControllerIp.text == '') ||
          (textControllerReadTime.text == '') ||
          (textControllerPort.text == '') ||
          (textControllerName.text == '')) {
        errorConfiguration = true;
        setState(() {
          errorConfiguration;
        });
        alertMessage.showError(context, translator["newDevices"]["errorDeviceMissing"]);
        return false;
      }
      if (topicsList.isEmpty) {
        errorTopics = true;
        setState(() {
          errorTopics;
        });
        alertMessage.showError(context, translator["newDevices"]["errorTopicMissing"]);
        return false;
      }

      List<Map<String, dynamic>> extraDataList = [];
      try {
        extraDataList
            .add(fetchDynamicControllerText(extraDataControllerMap, extraDataControllerMap.entries));
      } catch (e) {
        errorExtraData = true;
        extraDataControllerMap = e;
        setState(() {
          errorExtraData;
          extraDataControllerMap;
        });
        rebuildAllChildren(context);
        // errorKey = e.toString();
        // setState(() {
        //   errorKey;
        // });
        // loadExtraFields();
        alertMessage.showError(context, translator["newDevices"]["errorExtraDataMissing"]);
        return false;
      }

      JsonData jsonData = JsonData(
          address: textControllerIp.text,
          port: textControllerPort.text,
          name: textControllerName.text,
          protocolConnection: selectedProtocol!,
          readingTime: textControllerReadTime.text,
          topics: topicsList,
          data: extraDataList);

      Map<String, dynamic> json = jsonData.toJson();

      if (id == null) {
        Response response = await authService.saveDevice(jsonEncode(json));
        code = response.statusCode;
      } else {
        code = await authService.updateDevice(id.toString(), jsonEncode(json));
      }

      if (code == 200) {
        if (context.mounted) {
          alertMessage.showSuccess(context, translator["newDevices"]["successSave"]);
        }
        return true;
      } else {
        if (context.mounted) {
          alertMessage.showError(context, translator["newDevices"]["errorSaveError"]);
        }
        return false;
      }
    }

    confirmCallBack() => {
      if (onTap != null) {
        onTap()
      },
    };

    BlurryDialog  alert = BlurryDialog(
        title: translator["newDevices"]["alertTitle"],
        content: translator["newDevices"]["alertContent"],
        ButtonConfirmText: translator["newDevices"]["alertConfirmButton"],
        ButtonCancelText: translator["newDevices"]["alertCancelButton"],
        confirmCallBack: confirmCallBack
    );

    return Row(
      mainAxisAlignment: MainAxisAlignment.spaceBetween,
      children: [
        Container(
          height: 100,
          alignment: Alignment.centerLeft,
          child: Row(
            children: [
              CustomRoundedButton(
                textName: translator["newDevices"]["backButton"],
                height: 35,
                width: 120,
                fontSize: 14,
                isSelected: true,
                textColorActived: CustomColors.background700,
                textColorInactive: CustomColors.background700,
                splashColor: CustomColors.primary500,
                backgroundColorActived: CustomColors.primary500,
                backgroundColorInactive: CustomColors.primary500,
                borderRadiusValue: 20,
                onTap: () {
                  showDialog(
                    context: context,
                    builder: (BuildContext context) {
                      return alert;
                    },
                  );
                  // if (onTap != null) {
                  //   return onTap();
                  // }
                },
              ),
              Padding(
                padding: const EdgeInsets.only(left: 38.0),
                child: Column(
                  mainAxisAlignment: MainAxisAlignment.center,
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Text(
                      translator["newDevices"]["titleNewDevice"],
                      style: const TextStyle(
                          color: CustomColors.primary500,
                          fontSize: 24,
                          fontWeight: FontWeight.bold),
                    ),
                    Padding(
                        padding: const EdgeInsets.only(top: 8.0),
                        child: Text(
                          translator["newDevices"]["subtitleNewDevice"],
                          style: const TextStyle(
                              color: CustomColors.neutral700, fontSize: 14),
                        ))
                  ],
                ),
              ),
            ],
          ),
        ),
        Padding(
          padding: const EdgeInsets.only(right: 16.0),
          child: CustomRoundedButton(
            textName: translator["newDevices"]["saveButton"],
            height: 35,
            width: 120,
            fontSize: 14,
            isSelected: true,
            textColorActived: CustomColors.background700,
            textColorInactive: CustomColors.background700,
            splashColor: selectedProtocol == null
                ? CustomColors.primary500.withOpacity(0.5)
                : CustomColors.primary500,
            backgroundColorActived: selectedProtocol == null
                ? CustomColors.primary500.withOpacity(0.5)
                : CustomColors.primary500,
            backgroundColorInactive: selectedProtocol == null
                ? CustomColors.primary500.withOpacity(0.5)
                : CustomColors.primary500,
            borderRadiusValue: 20,
            onTap: () async {
              bool result = await save(widget.dataEdit != null ? widget.dataEdit['id'] : null);
              if (onTap != null && result) {
                onTap();
              }
            },
          ),
        ),
      ],
    );
  }

  Widget backButton(double heightHead) {
    return Consumer<LanguageProvider>(
      builder: (context, language, child) {
        final translator = language.getDataLanguage(language.currentLanguage);
        return Container(
          height: heightHead,
          alignment: Alignment.centerLeft,
          child: Row(
            children: [
              CustomRoundedButton(
                textName: translator["newDevices"]["backButton"],
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
                  setShowExtraData(false);
                  // for(dynamic entry in extraDataControllerMap.entries) {
                  //   resetColumnsSelected(entry.value);
                  // }
                },
              ),
              Padding(
                padding: const EdgeInsets.only(left: 16.0),
                child: Column(
                  mainAxisAlignment: MainAxisAlignment.center,
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Text(
                      translator["newDevices"]["titleExtraData"],
                      style: const TextStyle(
                          color: CustomColors.primaryColorApp,
                          fontSize: 16,
                          fontWeight: FontWeight.bold),
                    ),
                    Text(
                      translator["newDevices"]["subtitleExtraData"],
                      style: const TextStyle(
                          color: CustomColors.primaryColorApp, fontSize: 13),
                    )
                  ],
                ),
              )
            ],
          ),
        );
      }
    );
  }

  @override
  Widget build(BuildContext context) {
    const double heightHead = 68;
    return Container(
      padding: const EdgeInsets.only(left: 16.0, top: 8.0),
      child: showExtraData
          ? LayoutBuilder(
              builder: (context, size) {
                return Padding(
                  padding: const EdgeInsets.only(left: 16.0, top: 8.0),
                  child: Column(
                    children: [
                      backButton(heightHead - 8),
                      const SizedBox(height: heightHead / 2),
                      SizedBox(
                        height: size.maxHeight - (heightHead + heightHead / 2),
                        width: size.maxWidth,
                        child: SingleChildScrollView(
                          scrollDirection: Axis.horizontal,
                          child: Row(children: extraDataRowList),
                        ),
                      ),
                    ],
                  ),
                );
              },
            )
          : SingleChildScrollView(
              scrollDirection: Axis.vertical,
              child: Column(
                mainAxisAlignment: MainAxisAlignment.start,
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  saveButton(widget.onTap, context),
                  const SizedBox(height: 20),
                  SingleChildScrollView(
                    scrollDirection: Axis.horizontal,
                    child: Row(
                      mainAxisAlignment: MainAxisAlignment.start,
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: [
                        CardProtocol(
                          protocols: protocolList,
                          value: selectedProtocol,
                          onChanged: (value) {
                            selectedProtocol = value;
                            setState(() {
                              selectedProtocol;
                            });
                            loadExtraFields();
                          },
                          isErrored: errorProtocol,
                        ),
                        const SizedBox(width: 8.0),
                        CardTextFieldProtocol(
                          textControllerIp: textControllerIp,
                          textControllerName: textControllerName,
                          textControllerPort: textControllerPort,
                          textControllerReadTime: textControllerReadTime,
                          isErrored: errorConfiguration,
                        ),
                        const SizedBox(width: 8.0),
                        CardTopics(
                            topicsList: topicsList,
                            onChanged: (value) {
                              topicsList = value ?? [];
                              setState(() {
                                topicsList;
                              });
                            },
                            isErrored: errorTopics,
                        ),
                        const SizedBox(width: 8.0),
                        extraDataObj != null
                            ? CardExtraData(
                                extraData: widget.extraData,
                                extraDataObj: extraDataObj,
                                setShowExtraData: setShowExtraData,
                                setExtraDataRowList: setExtraDataRowList,
                                getExtraDataRowList: getExtraDataRowList,
                                isErrored: errorExtraData,
                              )
                            : Container(),
                      ],
                    ),
                  ),
                ],
              ),
            ),
    );
  }
}
