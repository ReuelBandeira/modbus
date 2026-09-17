import 'package:flutter/material.dart';
import 'package:isi_4_0/colors/color.dart';
import 'package:isi_4_0/model/cards_model.dart';
import 'package:isi_4_0/providers/lang.dart';
import 'package:isi_4_0/repository/db_api.dart';
import 'package:isi_4_0/viewmodel/device_view_model.dart';
import 'package:isi_4_0/views/devices/components/custom_cards.dart';
import 'package:isi_4_0/views/devices/components/custom_head_list_devices.dart';
import 'package:isi_4_0/views/devices/components/custom_navigator_page.dart';
import 'package:isi_4_0/views/devices/components/custom_row_list_device.dart';
import 'package:isi_4_0/views/devices/new_devices.dart';
import 'package:isi_4_0/widgets/custom_rounded_button.dart';
import 'package:provider/provider.dart';

class Devices extends StatefulWidget {
  const Devices({super.key, this.onTap});
  final Function()? onTap;
  @override
  State<Devices> createState() => _DevicesState();
}

class _DevicesState extends State<Devices> {
  bool pageEdit = false;
  dynamic dataDevice;
  int page = 1;
  int items = 10;
  int totalPages = 0;
  int totalItems = 0;

  @override
  void initState() {
    setState(() {
      page;
      items;
      totalPages;
      totalItems;
    });
    super.initState();
  }

  @override
  Widget build(BuildContext context) {
    Future<void> choose(
        String select, device, Map<String, dynamic> lang) async {
      switch (select) {
        case 'edit':
          setState(() {
            pageEdit = true;
            dataDevice = device;
          });
          break;
        case 'remove':
          if (await dialogRemove(context, device, lang)) {
            setState(() {
              pageEdit = false;
            });
          }
          break;
      }
    }

    void receiveDataFromChild(dynamic data) {
      setState(() {});
    }

    return pageEdit
        ? NewDevices(
            dataCallback: receiveDataFromChild,
            dataEdit: dataDevice,
            onTap: (() {
              setState(() {
                pageEdit = false;
              });
            }),
          )
        : Consumer2<DeviceViewModel, LanguageProvider>(
            builder: (context, devices, language, child) {
              final lang = language.getDataLanguage(language.currentLanguage);
              return Padding(
                padding: const EdgeInsets.only(top: 12.0, left: 12.0),
                child: SingleChildScrollView(
                  scrollDirection: Axis.vertical,
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      FutureBuilder<List<InfoCards>>(
                          future: devices.getDataCards(language),
                          builder: (context, snapshotCards) {
                            if (snapshotCards.connectionState ==
                                ConnectionState.waiting) {
                              return const SizedBox(
                                  height: 120,
                                  child: Center(
                                      child: Column(
                                        mainAxisAlignment: MainAxisAlignment.center,
                                        crossAxisAlignment: CrossAxisAlignment.center,
                                        children: [
                                          CircularProgressIndicator(),
                                        ],
                                      )
                                  )
                              );
                            } else if (snapshotCards.hasError) {
                              return SizedBox(
                                height: 120,
                                child: Center(
                                    child: Text(
                                      '${lang['devices']['error']}: ${snapshotCards.error}',
                                      style: const TextStyle(fontSize: 20),
                                    )
                                ),
                              );
                            } else if (!snapshotCards.hasData ||
                                snapshotCards.data!.isEmpty) {
                              return SizedBox(
                                height: 120,
                                child: Center(
                                    child: Text(
                                      lang['devices']['noCardsData'],
                                      style: const TextStyle(fontSize: 20),
                                    )
                                ),
                              );
                            } else {
                              List<InfoCards> infoCards = snapshotCards.data!;
                              return SingleChildScrollView(
                                scrollDirection: Axis.horizontal,
                                child: Row(
                                    children: infoCards.map((infoCard) {
                                  return CustomCards(
                                    info: infoCard,
                                  );
                                }).toList()),
                              );
                            }
                          }),
                      FutureBuilder<Map<String, dynamic>>(
                        future: devices.getDataDevices(page, items),
                        builder: (context, snapshotDevices) {
                          if (snapshotDevices.connectionState ==
                              ConnectionState.waiting) {
                            return const SizedBox(
                                height: 300,
                                child: Center(
                                    child: Column(
                                      mainAxisAlignment: MainAxisAlignment.center,
                                      crossAxisAlignment: CrossAxisAlignment.center,
                                      children: [
                                        CircularProgressIndicator(),
                                      ],
                                    )
                                )
                            );
                          } else if (snapshotDevices.hasError) {
                            return SizedBox(
                              height: 300,
                              child: Center(
                                  child: Text(
                                    '${lang['devices']['error']}: ${snapshotDevices.error}',
                                    style: const TextStyle(fontSize: 25),
                                  )
                              ),
                            );
                          } else if (!snapshotDevices.hasData ||
                              snapshotDevices.data!.isEmpty) {
                            return Column(
                                crossAxisAlignment: CrossAxisAlignment.start,
                                children: [
                                  CustomRoundedButton(
                                    textName: lang['devices']['newDevice'],
                                    height: 35,
                                    width: 120,
                                    fontSize: 14,
                                    isSelected: true,
                                    textColorActived: CustomColors.background700,
                                    textColorInactive: CustomColors.background700,
                                    splashColor: CustomColors.primary500,
                                    backgroundColorActived:
                                    CustomColors.primary500,
                                    backgroundColorInactive:
                                    CustomColors.primary500,
                                    borderRadiusValue: 20,
                                    onTap: widget.onTap,
                                  ),
                                  SizedBox(
                                    height: 300,
                                    child: Center(
                                        child: Text(
                                          lang['devices']['noData'],
                                          style: const TextStyle(fontSize: 25),
                                        )
                                    ),
                                  ),
                                ]
                            );
                          } else {
                            List<dynamic> getdevices = snapshotDevices.data!['items']!;
                            List<dynamic> deviceList = getdevices[0]['devices'];
                            totalPages = snapshotDevices.data!['pages'];
                            totalItems = snapshotDevices.data!['total'];

                            return Column(
                              crossAxisAlignment: CrossAxisAlignment.start,
                              children: [
                                CustomRoundedButton(
                                  textName: lang['devices']['newDevice'],
                                  height: 35,
                                  width: 120,
                                  fontSize: 14,
                                  marginTop: 20.0,
                                  isSelected: true,
                                  textColorActived: CustomColors.background700,
                                  textColorInactive: CustomColors.background700,
                                  splashColor: CustomColors.primary500,
                                  backgroundColorActived:
                                  CustomColors.primary500,
                                  backgroundColorInactive:
                                  CustomColors.primary500,
                                  borderRadiusValue: 20,
                                  onTap: widget.onTap,
                                ),
                                CustomHeadListDevices(language: language),
                                FittedBox(
                                  fit: BoxFit.fitHeight,
                                  child: Column(
                                    children: deviceList.map((device) {
                                      return CustomRowListDevice(
                                        device: device,
                                        buttonsName: lang,
                                        hint: lang['devices']['listTopics'],
                                        tooltipsMessage: lang['tooltips']
                                            ['topicsMenu'],
                                        onSelected: (select) {
                                          choose(select, device, lang);
                                        },
                                      );
                                    }).toList(),
                                  ),
                                ),
                                CustomNavigatorPage(
                                  language: language,
                                  page: page,
                                  items: items,
                                  totalPages: totalPages,
                                  totalItems: totalItems,
                                  onItemsChange: (int newItems) {
                                    items = newItems;
                                    page = 1;
                                    setState(() {
                                      items;
                                      page;
                                    });
                                  },
                                  onPageChange: (int newPage) {
                                    page = newPage;
                                    setState(() {
                                      page;
                                    });
                                  },
                                ),
                              ],
                            );
                          }
                        },
                      ),
                    ],
                  ),
                ),
              );
            },
          );
  }
}

Future<bool> dialogRemove(
    BuildContext context, device, Map<String, dynamic> lang) async {
  void close(bool state) {
    Navigator.of(context).pop(state);
  }

  return await showDialog(
    context: context,
    barrierDismissible: false,
    builder: (BuildContext context) {
      AuthService authService = AuthService();
      return Dialog(
        child: SizedBox(
          height: 200,
          width: 400,
          child: Column(
            children: [
              Container(
                  padding: const EdgeInsets.only(left: 16, top: 16, bottom: 16),
                  alignment: Alignment.topLeft,
                  child: Text(
                    lang['devices']['removeDevice'],
                    style: const TextStyle(
                        color: CustomColors.primaryColorApp,
                        fontWeight: FontWeight.bold,
                        fontSize: 18),
                  )),
              Padding(
                padding: const EdgeInsets.all(8.0),
                child: Text('${lang['devices']['alertDevice']}?'),
              ),
              Padding(
                padding: const EdgeInsets.all(8.0),
                child: Text(
                  device['name'],
                  style: const TextStyle(
                      color: CustomColors.error500,
                      fontWeight: FontWeight.bold),
                ),
              ),
              Padding(
                padding: const EdgeInsets.all(8.0),
                child: Row(
                  mainAxisAlignment: MainAxisAlignment.spaceEvenly,
                  children: [
                    TextButton(
                        onPressed: (() {
                          close(false);
                        }),
                        child: Text(lang['devices']['cancel'])),
                    CustomRoundedButton(
                      textName: lang['devices']['remove'],
                      height: 30,
                      width: 80,
                      fontSize: 12,
                      isSelected: true,
                      textColorActived: CustomColors.background700,
                      textColorInactive: CustomColors.background700,
                      splashColor: CustomColors.primary500,
                      backgroundColorActived: CustomColors.primary500,
                      backgroundColorInactive: CustomColors.primary500,
                      borderRadiusValue: 20,
                      onTap: (() async {
                        bool response = await authService
                            .deleteDevice(device['id'].toString());
                        close(response);
                      }),
                    )
                  ],
                ),
              )
            ],
          ),
        ),
      );
    },
  );
}
