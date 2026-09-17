import 'package:flutter/material.dart';
import 'package:isi_4_0/colors/color.dart';
import 'package:isi_4_0/widgets/custom_device_text_field.dart';
import 'package:isi_4_0/widgets/custom_rounded_button.dart';
import 'package:provider/provider.dart';
import '../../../providers/lang.dart';
import '../../../utils/alert_message.dart';
import 'custom_buttom_card_extra.dart';

class CardExtraData extends StatelessWidget {
  const CardExtraData(
      {super.key, this.loadFields, this.extraData, this.extraDataObj,
        required this.setShowExtraData, required this.setExtraDataRowList,
        required this.getExtraDataRowList, required this.isErrored});
  final Function()? loadFields;
  final Function()? extraData;
  final Widget? extraDataObj;
  final Function(bool) setShowExtraData;
  final Function(List<Widget>) setExtraDataRowList;
  final Function() getExtraDataRowList;
  final bool isErrored;

  Widget buildDynamicField(
      dynamic extraField,
      dynamic extraDataObj,
      dynamic extraDataObjValue,
      dynamic extraDataObjKey,
      dynamic controllerMap,
      void Function(dynamic) dataCallback,
      [bool isRoot = false]) {
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
          text: extraDataObjValue == null ? '' : extraDataObjValue.toString());
      controllerMap[extraDataObjKey.toString()] = {
        'controller': controller,
        'type': extraField.value['type'],
        'name': extraField.value['name'],
        'name_pt': extraField.value['name_pt'],
        'name_es': extraField.value['name_es'],
        'min': extraField.value['min'],
        'max': extraField.value['max'],
        'required': ((extraField.value['requirements'].length > 0) &&
            extraField.value['requirements'].contains('nonNull')),
        'isErrored': false,
      };
      // controllerList.add(controller);
      return Consumer<LanguageProvider>(
          builder: (context, language, child) {
            String fieldName = extraField.value['name'];
            if (language.currentLanguage != 'en') {
              if (extraField.value['name_${language.currentLanguage}'] != null) {
                fieldName = extraField.value['name_${language.currentLanguage}'];
              }
            }
            return Column(children: [
                const SizedBox(height: 12),
                Container(
                    // height: isRoot ? 50 : 38,
                    padding: isRoot ? const EdgeInsets.only(bottom: 12.0) : null,
                    decoration: isRoot
                        ? const BoxDecoration(
                            border: Border(
                                bottom: BorderSide(color: Colors.black26, width: 1)))
                        : null,
                    child: CustomDeviceTextField(
                      controller: controller,
                      min: extraField.value['min'],
                      max: extraField.value['max'],
                      isRequired: ((extraField.value['requirements'].length > 0) &&
                          extraField.value['requirements'].contains('nonNull')),
                      digitsOnly: ((extraField.value['type'] == "int") ||
                          (extraField.value['type'] == "float")),
                      labelText: fieldName,
                    )),
              ]
            );
          }
      );
    } else {
      if (extraField.value['type'] == "array") {
        dynamic mainObj = {};
        for (var object in extraField.value.entries) {
          if (object.key == 'contains') {
            mainObj = object;
          }
        }
        controllerMap[extraDataObjKey.toString()] = {
          'controller': [],
          'type': extraField.value['type'],
          'data': [],
          'hover': [],
          'select': [],
          'min': extraField.value['min'],
          'max': extraField.value['max'],
          'required': ((extraField.value['requirements'].length > 0) &&
              extraField.value['requirements'].contains('nonNull')),
          'isErrored': false,
        };

        List<Widget> createChildrenList() {
          List<Widget> list = List.empty(growable: true);

          if (extraDataObjValue != null) {
            for (var i = 0; i < extraDataObjValue.length; i++) {
              if (controllerMap[extraDataObjKey.toString()]['controller'].length-1 < i) {
                controllerMap[extraDataObjKey.toString()]['controller'].add({});
                controllerMap[extraDataObjKey.toString()]['hover'].add(false);
                controllerMap[extraDataObjKey.toString()]['select'].add(false);
              }
              Widget listItem = buildDynamicField(
                mainObj,
                extraDataObjValue,
                extraDataObjValue[i],
                extraDataObjKey,
                controllerMap[extraDataObjKey.toString()]['controller'][i],
                dataCallback,
              );

              // Widget listItem = customColumn(
              //     child: buildDynamicField(
              //           mainObj,
              //           extraDataObjValue,
              //           extraDataObjValue[i],
              //           extraDataObjKey,
              //           controllerMap[extraDataObjKey.toString()]['controller'][i],
              //           dataCallback,
              //       ),
              //     data: extraDataObjValue[i],
              //     controller: controllerMap[extraDataObjKey.toString()]['controller'][i],
              //     columnIndex: i,
              //     hover: controllerMap[extraDataObjKey.toString()]['hover'],
              //     select: controllerMap[extraDataObjKey.toString()]['select'],
              //     buttonNewVisible: true,
              //     columnVisible: true,
              //     columnName: extraField.value['name'].toString(),
              //     icon: true
              // );

              // Widget listItem = Column(
              //   children: [
              //     Card(
              //       child: buildDynamicField(
              //           mainObj,
              //           extraDataObjValue,
              //           extraDataObjValue[i],
              //           extraDataObjKey,
              //           controllerMap[extraDataObjKey.toString()]['controller'][i],
              //           dataCallback,
              //           extraDataRowList,
              //       )
              //     ),
              //     ElevatedButton(
              //       style: ElevatedButton.styleFrom(
              //         primary: Colors.lightBlue,
              //         padding: const EdgeInsets.all(12),
              //         textStyle: const TextStyle(fontSize: 22),
              //       ),
              //       onPressed: () {
              //         extraDataObjValue.removeAt(i);
              //         // updatePage();
              //       },
              //       child: const Text('Remove'),
              //     )
              //   ],
              // );

              list.add(listItem);
            }
          }

          return list;
        }

        List<Widget> childrenList = createChildrenList();


        // // DEBUG
        // controllerMap[extraDataObjKey.toString()]['controller'].add({});
        // controllerMap[extraDataObjKey.toString()]['hover'].add(false);
        // controllerMap[extraDataObjKey.toString()]['select'].add(false);
        // Widget listItem = customColumn(
        //     data: [
        //       {
        //         "address": 40964,
        //         "name": "Y0.4"
        //       }
        //     ],
        //     child: buildDynamicField(
        //       mainObj,
        //       [
        //         {
        //           "address": 40964,
        //           "name": "Y0.4"
        //         }
        //       ],
        //       {
        //         "address": 40964,
        //         "name": "Y0.4"
        //       },
        //       extraDataObjKey,
        //       controllerMap[extraDataObjKey.toString()]['controller'][0],
        //       dataCallback,
        //       showExtraData,
        //     ),
        //     controller: controllerMap[extraDataObjKey.toString()]['controller'][0],
        //     columnIndex: 0,
        //     hover: controllerMap[extraDataObjKey.toString()]['hover'],
        //     select: controllerMap[extraDataObjKey.toString()]['select'],
        //     buttonNewVisible: true,
        //     columnVisible: true,
        //     columnName: extraField.value['name'].toString(),
        //     icon: true
        // );
        //
        //
        //
        //
        //
        // // Widget listItem = Column(
        // //   children: [
        // //     Card(
        // //         child: buildDynamicField(
        // //           mainObj,
        // //           extraDataObjValue,
        // //           null,
        // //           extraDataObjKey,
        // //           controllerMap[extraDataObjKey.toString()]['controller'][0],
        // //           dataCallback,
        // //           extraDataRowList,
        // //         )
        // //     ),
        // //     ElevatedButton(
        // //       style: ElevatedButton.styleFrom(
        // //         primary: Colors.lightBlue,
        // //         padding: const EdgeInsets.all(12),
        // //         textStyle: const TextStyle(fontSize: 22),
        // //       ),
        // //       onPressed: () {
        // //         extraDataObjValue.removeAt(0);
        // //         // updatePage();
        // //       },
        // //       child: const Text('Remove'),
        // //     )
        // //   ],
        // // );
        //
        //
        //
        // dataList.add(listItem);


        if (childrenList.isNotEmpty) {
          List<bool> hover = List.filled(childrenList.length, false, growable: true);
          List<bool> select = List.filled(childrenList.length, false, growable: true);
          controllerMap[extraDataObjKey.toString()]['hover'] = hover;
          controllerMap[extraDataObjKey.toString()]['select'] = select;
        }

        controllerMap[extraDataObjKey.toString()]['data'] = childrenList;

        if (isRoot) {
          return Container(
            decoration: controllerMap[extraDataObjKey.toString()]['isErrored'] ? BoxDecoration(
                border: Border.all(color: CustomColors.error600Fade, width: 1),
                borderRadius: const BorderRadius.all(
                    Radius.circular(5.0)
                ),
            ) : null,
            child: Consumer<LanguageProvider>(
              builder: (context, language, child) {
                final translator = language.getDataLanguage(language.currentLanguage);

                String fieldName = extraField.value['name'];
                if (language.currentLanguage != 'en') {
                  if (extraField.value['name_${language.currentLanguage}'] != null) {
                    fieldName = extraField.value['name_${language.currentLanguage}'];
                  }
                }

                Widget starterColumn = customColumn(
                    children: childrenList,
                    data: extraDataObjValue,
                    removeExtraData: (int index) {
                      extraDataObjValue.removeAt(index);
                      controllerMap[extraDataObjKey.toString()]['controller'].removeAt(index);
                      // controllerMap[extraDataObjKey.toString()]['hover'].removeAt(index);
                      // controllerMap[extraDataObjKey.toString()]['select'].removeAt(index);
                    },
                    addExtraData: (dynamic data) {
                      if (extraDataObjValue == null) {
                        extraDataObj[extraDataObjKey] = [];
                        extraDataObjValue = extraDataObj[extraDataObjKey];
                      }
                      extraDataObjValue = data;

                      controllerMap[extraDataObjKey.toString()]['controller'].add({});
                      controllerMap[extraDataObjKey.toString()]['hover'].add(false);
                      controllerMap[extraDataObjKey.toString()]['select'].add(false);

                      // print(extraDataObjValue);
                      // print(controllerMap[extraDataObjKey.toString()]['controller']);

                      Widget listItem = buildDynamicField(
                        mainObj,
                        extraDataObjValue,
                        extraDataObjValue[extraDataObjValue.length - 1],
                        extraDataObjKey,
                        controllerMap[extraDataObjKey.toString()]['controller'][extraDataObjValue.length - 1],
                        dataCallback,
                      );

                      // controllerMap[extraDataObjKey.toString()]['data'].add(listItem);

                      return listItem;
                      // updatePage();
                    },
                    columnIndex: 0,
                    hover: controllerMap[extraDataObjKey.toString()]['hover'],
                    select: controllerMap[extraDataObjKey.toString()]['select'],
                    buttonNewVisible: true,
                    columnVisible: true,
                    columnName: fieldName.toString(),
                    icon: true,
                    min: extraField.value['min'],
                    max: extraField.value['max']
                );

                return InkWell(
                  onTap: () {
                    List<Widget> extraDataRowList = [
                      starterColumn
                    ];
                    setExtraDataRowList(extraDataRowList);
                    setShowExtraData(true);
                  },
                  child: Column(children: [
                    const SizedBox(height: 12),
                    Container(
                      height: 50,
                      padding: const EdgeInsets.only(bottom: 12.0),
                      decoration: const BoxDecoration(
                          border: Border(
                              bottom: BorderSide(color: Colors.black26, width: 1))),
                      child: Row(
                        mainAxisAlignment: MainAxisAlignment.spaceBetween,
                        children: [
                          Text(
                            fieldName.toString(),
                            style: const TextStyle(color: Colors.grey),
                          ),
                          CustomRoundedButton(
                            textName: translator["newDevices"]["buttonSee"],
                            height: 32,
                            width: 50,
                            fontSize: 13,
                            isSelected: true,
                            textColorActived: CustomColors.background700,
                            textColorInactive: CustomColors.background700,
                            splashColor: CustomColors.primary500,
                            backgroundColorActived: CustomColors.primary500,
                            backgroundColorInactive: CustomColors.primary500,
                            borderRadiusValue: 20,
                            onTap: () {
                              List<Widget> extraDataRowList = [
                                starterColumn
                              ];
                              setExtraDataRowList(extraDataRowList);
                              setShowExtraData(true);
                            },
                          )
                        ],
                      ),
                    ),

                    //   Text(extraField.key.toString()),
                    //   SingleChildScrollView(
                    //       scrollDirection: Axis.vertical,
                    //       child: Card(
                    //           child: Column(
                    //             children: [
                    //               ...data(),
                    //               ElevatedButton(
                    //                 style: ElevatedButton.styleFrom(
                    //                   primary: Colors.lightBlue,
                    //                   padding: const EdgeInsets.all(12),
                    //                   textStyle: const TextStyle(fontSize: 22),
                    //                 ),
                    //                 onPressed: () {
                    //                   if (extraDataObjValue == null) {
                    //                     extraDataObj[extraDataObjKey] = [];
                    //                     extraDataObjValue = extraDataObj[extraDataObjKey];
                    //                   }
                    //                   extraDataObjValue.add({});
                    //                   // updatePage();
                    //                 },
                    //                 child: const Text('Add'),
                    //               )
                    //             ],
                    //           )))
                  ]),
                );
              }
            )
          );
        } else {
          //TODO: array as Child of array
          // return customColumn(
          //     extraDataObjValue: controllerMap[extraDataObjKey.toString()]['data'],
          //     extraField: mainObj,
          //     controller: controllerMap[extraDataObjKey.toString()],
          //     hover: controllerMap[extraDataObjKey.toString()]['hover'],
          //     select: controllerMap[extraDataObjKey.toString()]['select'],
          //     buttonNewVisible: true,
          //     collumnVisible: true,
          //     columnName: extraField.value['name'].toString(),
          //     icon: true
          // );
        }
      } else {
        // object
        controllerMap[extraDataObjKey.toString()] = {
          'controller': {},
          'type': extraField.value['type'],
          'name': extraField.value['name'],
          'name_pt': extraField.value['name_pt'],
          'name_es': extraField.value['name_es'],
          'min': extraField.value['min'],
          'max': extraField.value['max'],
          'required': ((extraField.value['requirements'].length > 0) &&
              extraField.value['requirements'].contains('nonNull')),
          'isErrored': false,
        };

        List<Widget> data() {
          List<Widget> list = List.empty(growable: true);
          for (var object in extraField.value.entries) {
            if (object.key == 'contains') {
              for (var objectField in object.value.entries) {
                dynamic objValue;
                if (extraDataObjValue != null) {
                  objValue = extraDataObjValue[objectField.key];
                }
                list.add(buildDynamicField(
                    objectField,
                    extraDataObjValue,
                    objValue,
                    objectField.key,
                    controllerMap[extraDataObjKey.toString()]['controller'],
                    dataCallback,
                  )
                );
              }
            }
          }
          return list;
        }

        if (isRoot) {
          return Consumer<LanguageProvider>(
              builder: (context, language, child) {
                final translator = language.getDataLanguage(language.currentLanguage);
                String? fieldName = extraField.value['name'];
                if (language.currentLanguage != 'en') {
                  if (extraField.value['name_${language.currentLanguage}'] != null) {
                    fieldName = extraField.value['name_${language.currentLanguage}'];
                  }
                }
                return Container(
                    padding: const EdgeInsets.only(bottom: 12.0),
                    decoration: const BoxDecoration(
                        border: Border(
                            bottom: BorderSide(color: Colors.black26, width: 1))),
                    child: Column(
                      crossAxisAlignment: CrossAxisAlignment.stretch,
                      children: [
                        const SizedBox(height: 12),
                        Text(
                          fieldName != null ?
                          fieldName.toString() :
                          translator["newDevices"]["labelItem"],
                          style: const TextStyle(color: Colors.grey),
                        ),
                        ...data()
                      ],
                    ));
              }
          );
          // return InkWell(
          //   onTap: extraData, // TODO: Call next screen passing data()
          //   child: Column(children: [
          //     const SizedBox(height: 12),
          //     Container(
          //       height: 50,
          //       padding: const EdgeInsets.only(bottom: 12.0),
          //       decoration: const BoxDecoration(
          //           border: Border(
          //               bottom:
          //               BorderSide(color: Colors.black26, width: 1))),
          //       child: Row(
          //         mainAxisAlignment: MainAxisAlignment.spaceBetween,
          //         children: [
          //           Text(
          //             extraField.key.toString(),
          //             style: const TextStyle(color: Colors.grey),
          //           ),
          //           CustomRoundedButton(
          //             textName: 'SEE',
          //             height: 32,
          //             width: 50,
          //             fontSize: 13,
          //             isSelected: true,
          //             textColorActived: CustomColors.background700,
          //             textColorInactive: CustomColors.background700,
          //             splashColor: CustomColors.primary500,
          //             backgroundColorActived: CustomColors.primary500,
          //             backgroundColorInactive: CustomColors.primary500,
          //             borderRadiusValue: 20,
          //             onTap: extraData, // TODO: Call next screen passing data()
          //           )
          //         ],
          //       ),
          //     ),
          //   ]),
          // );
        } else {
          return Container(
              padding: const EdgeInsets.only(bottom: 12.0),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.stretch,
                children: [
                  // const SizedBox(height: 12),
                  ...data()
                ],
              ));
        }
      }

      // All else failed, unexpected situation, show nothing.
      return Container();
    }
  }

  Widget customColumn(
      {required List<dynamic>? data,
        required List<dynamic>? children,
        required int columnIndex,
        required List<dynamic>? hover,
        required List<dynamic>? select,
        Function(int)? removeExtraData,
        Function(dynamic)? addExtraData,
        bool? buttonNewVisible,
        bool? columnVisible,
        String? columnName,
        bool? icon,
        bool showChild = false,
        int? childIndex,
        int? min,
        int? max,
      }) {
    const double spaces = 48.0;
    const double width = 300;
    AlertMessage alertMessage = AlertMessage();

    BuildContext? dataContext;

    return Consumer<LanguageProvider>(
        builder: (context, language, child) {
          final translator = language.getDataLanguage(language.currentLanguage);
          return Container(
            alignment: Alignment.topCenter,
            margin: const EdgeInsets.only(right: spaces / 2),
            decoration: const BoxDecoration(
                border: Border(
                    right: BorderSide(color: CustomColors.neutral400, width: 1.0))),
            child: Stack(
              children: [
                Align(
                  alignment: Alignment.topCenter,
                  child: Container(
                    height: spaces,
                    width: width,
                    alignment: Alignment.topCenter,
                    child: Text(
                      columnName ?? translator["newDevices"]["errorMissingColumn"],
                      style: const TextStyle(color: CustomColors.primaryColorApp),
                    ),
                  ),
                ),
                Align(
                  alignment: Alignment.center,
                  child: Visibility(
                    visible: columnVisible ?? false,
                    child: Container(
                      padding: const EdgeInsets.only(
                          top: spaces, bottom: spaces * 2, right: spaces / 2),
                      alignment: Alignment.topCenter,
                      child: SingleChildScrollView(
                        scrollDirection: Axis.vertical,
                        child: StatefulBuilder(builder: (context, buttonState) {
                          dataContext = context;
                          return showChild && (childIndex != null) ?
                          SizedBox(
                              width: width,
                              child: children?[childIndex]
                          ) : (data != null) && data!.isNotEmpty ?
                          Column(
                              children: List.generate(data!.length, (index) {
                                return InkWell(
                                  child: CustomButtomCardExtra(
                                    height: spaces + 8,
                                    width: width,
                                    field: translator["newDevices"]["labelItem"] + ' $index',
                                    hover: hover?[index],
                                    select: select?[index],
                                    showIcon: icon ?? false,
                                    onPressedDelete: () {
                                      if (min != null) {
                                        if (data!.length - 1 < min) {
                                          alertMessage.showError(
                                            context,
                                            translator["newDevices"]["errorMinElements"]
                                              + min.toString()
                                              + translator["newDevices"]["labelElements"]
                                          );
                                          return;
                                        }
                                      }
                                      if (select?[index] == true) {
                                        List<
                                            Widget> extraDataRowList = getExtraDataRowList();

                                        for (int i = 0; i + columnIndex + 1 <
                                            extraDataRowList.length; i++) {
                                          extraDataRowList[i + columnIndex + 1] =
                                              Container();
                                        }

                                        setExtraDataRowList(extraDataRowList);
                                      }
                                      // buttonState(() => data!.removeAt(index));
                                      children?.removeAt(index);
                                      hover?.removeAt(index);
                                      select?.removeAt(index);
                                      if (removeExtraData != null) {
                                        removeExtraData(index);
                                      }
                                      if (dataContext != null) {
                                        (dataContext as Element).markNeedsBuild();
                                      }
                                    },
                                  ),
                                  onTap: () {
                                    select =
                                        select?.map<bool>((v) => false).toList();
                                    buttonState(() =>
                                    select?[index] = !select?[index]);

                                    Widget column = customColumn(
                                      data: data,
                                      children: children,
                                      columnIndex: columnIndex+1,
                                      hover: hover,
                                      select: select,
                                      childIndex: index,
                                      removeExtraData: removeExtraData,
                                      addExtraData: addExtraData,
                                      buttonNewVisible: false,
                                      columnVisible: true,
                                      columnName: translator["newDevices"]["labelItem"] + ' $index',
                                      icon: false,
                                      showChild: true,
                                      min: min,
                                      max: max,
                                    );

                                    List<
                                        Widget> extraDataRowList = getExtraDataRowList();

                                    if (extraDataRowList.length <=
                                        columnIndex + 1) {
                                      extraDataRowList.add(
                                          column
                                      );
                                    } else {
                                      extraDataRowList[columnIndex + 1] = column;
                                      for (int i = 1; i + columnIndex + 1 <
                                          extraDataRowList.length; i++) {
                                        extraDataRowList[i + columnIndex + 1] =
                                            Container();
                                      }
                                    }

                                    setExtraDataRowList(extraDataRowList);
                                  },
                                  onHover: (value) {
                                    buttonState(() => hover?[index] = value);
                                  },
                                );
                              })
                          ) :
                          SizedBox(
                              width: width,
                              child: Text(translator["newDevices"]["errorEmptyList"],
                                  textAlign: TextAlign.center,
                                  style: const TextStyle(
                                      fontSize: 16,
                                      color: CustomColors.neutral700Faded,
                                      fontWeight: FontWeight.bold))
                          );
                        }),
                      ),
                    ),
                  ),
                ),
                Align(
                  alignment: Alignment.bottomCenter,
                  child: Container(
                    alignment: Alignment.center,
                    height: 88,
                    width: width,
                    child: Visibility(
                      visible: buttonNewVisible ?? false,
                      child: CustomRoundedButton(
                        textName: translator["newDevices"]["buttonNewItem"],
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
                          data ??= [];
                          hover ??= [];
                          select ??= [];
                          if (max != null) {
                            if (data!.length + 1 > max) {
                              alertMessage.showError(context,
                                  translator["newDevices"]["errorMaxElements"]
                                      + max.toString()
                                      + translator["newDevices"]["labelElements"]
                              );
                              return;
                            }
                          }
                          data?.add({});
                          hover?.add(false);
                          select?.add(false);
                          if (addExtraData != null) {
                            children?.add(addExtraData(data));
                          }
                          if (dataContext != null) {
                            (dataContext as Element).markNeedsBuild();
                          }
                        },
                      ),
                    ),
                  ),
                ),
              ],
            ),
          );
        }
    );
  }

  @override
  Widget build(BuildContext context) {
    return Consumer<LanguageProvider>(
        builder: (context, language, child) {
          final translator = language.getDataLanguage(language.currentLanguage);
          return SizedBox(
            height: 500,
            width: 350,
            child: Card(
              elevation: 2,
              color: CustomColors.background700,
              shape: RoundedRectangleBorder(
                borderRadius: BorderRadius.circular(2.0),
                side: BorderSide(
                  color: isErrored ? CustomColors.error600Fade : Colors.black12,
                ),
              ),
              child: Container(
                padding: const EdgeInsets.all(16),
                child: Column(
                  mainAxisAlignment: MainAxisAlignment.start,
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Container(
                        alignment: Alignment.bottomLeft,
                        padding: const EdgeInsets.only(bottom: 8.0),
                        child: Text(translator["newDevices"]["titleExtraData"],
                            style: const TextStyle(
                                fontSize: 16,
                                color: CustomColors.neutral800,
                                fontWeight: FontWeight.bold))),
                    Container(
                        alignment: Alignment.topLeft,
                        child: Text(translator["newDevices"]["subtitleExtraDataCard"],
                            style: const TextStyle(
                                fontSize: 12, color: CustomColors.neutral700))),
                    Expanded(
                      child: SingleChildScrollView(
                        scrollDirection: Axis.vertical,
                        child: Column(
                          children: [
                            const SizedBox(height: 20),
                            extraDataObj ?? Container(),
                          ],
                        ),
                      )
                    )
                  ],
                ),
              ),
            ),
          );
        }
    );
  }
}
