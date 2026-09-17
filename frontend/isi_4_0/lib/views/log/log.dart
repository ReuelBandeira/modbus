import 'package:flutter/material.dart';
import 'package:isi_4_0/colors/color.dart';
import 'package:isi_4_0/providers/lang.dart';
import 'package:isi_4_0/viewmodel/log_view_model.dart';
import 'package:isi_4_0/views/devices/components/custom_navigator_page.dart';
import 'package:isi_4_0/views/log/components/custom_head_list_files.dart';
import 'package:isi_4_0/views/log/components/custom_rows_files.dart';
import 'package:isi_4_0/widgets/custom_rounded_button.dart';
import 'package:provider/provider.dart';

import '../../utils/alert_message.dart';

class Log extends StatefulWidget {
  const Log({super.key});

  @override
  State<Log> createState() => _LogState();
}

class _LogState extends State<Log> {
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
    AlertMessage alertMessage = AlertMessage();

    return Consumer2<LogViewModel, LanguageProvider>(
      builder: (context, logs, language, child) {
        return Padding(
          padding: const EdgeInsets.only(top: 12.0, left: 12.0),
          child: SingleChildScrollView(
            scrollDirection: Axis.vertical,
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                FutureBuilder<Map<String, dynamic>>(
                  future: Future(() => logs.getLogList(page, items)),
                  builder: (context, files) {
                    if (files.connectionState == ConnectionState.waiting) {
                      return const SizedBox(
                          height: 420,
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
                    } else if (files.hasError) {
                      return SizedBox(
                        height: 200,
                        child: Center(
                            child: Text(
                              '${language.getDataLanguage(language.currentLanguage)
                              ['devices']['error']}: ${files.error}',
                              style: const TextStyle(fontSize: 25),
                            )
                        ),
                      );
                    } else if (!files.hasData) {
                      return SizedBox(
                        height: 200,
                        child: Center(
                            child: Text(
                              language.getDataLanguage(language.currentLanguage)
                              ['logs']['noLogs'],
                              style: const TextStyle(fontSize: 25),
                            )
                        ),
                      );
                    } else {
                      List<dynamic> logList = files.data!['items'];
                      totalPages = files.data!['pages'];
                      totalItems = files.data!['total'];
                      return Column(
                        crossAxisAlignment: CrossAxisAlignment.start,
                        children: [
                          CustomRoundedButton(
                            textName: language
                                .getDataLanguage(language.currentLanguage)['devices']
                            ['button_log_all_download'],
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
                              logs.downloadAllLogFile().
                              then((value) {
                                if (context.mounted) {
                                  alertMessage.showSuccess(context, language
                                      .getDataLanguage(language.currentLanguage)
                                  ['logs']['downloadSuccess']);
                                }
                              }).catchError((e){
                                if (context.mounted) {
                                  alertMessage.showError(context, '${language
                                      .getDataLanguage(language.currentLanguage)
                                  ['devices']['error']}: ${e.toString()}');
                                }
                              });
                            },
                          ),
                          CustomHeadListFiles(language: language),
                          Column(
                            children: logList.map((file) {
                              return CustomRowsFiles(
                                nameFile: file['name'],
                                dateFile: file['date'],
                                tooltipsMessage: language.getDataLanguage(
                                    language.currentLanguage)['logs']
                                ['downloadLog'],
                                onPressed: () {
                                  logs.downloadLogFile(file['id']).
                                  then((value) {
                                    if (context.mounted) {
                                      alertMessage.showSuccess(context, language
                                          .getDataLanguage(language.currentLanguage)
                                      ['logs']['downloadSuccess']);
                                    }
                                  }).catchError((e){
                                    if (context.mounted) {
                                      alertMessage.showError(context, '${language
                                          .getDataLanguage(language.currentLanguage)
                                      ['devices']['error']}: ${e.toString()}');
                                    }
                                  });
                                },
                              );
                            }).toList()),
                          CustomNavigatorPage(
                            language: language,
                            paddingRight: 10.0,
                            page: page,
                            items: items,
                            totalPages: totalPages,
                            totalItems: totalItems,
                            onItemsChange: (int newItems) {
                              items = newItems;
                              setState(() {
                                items;
                                page = 1;
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
