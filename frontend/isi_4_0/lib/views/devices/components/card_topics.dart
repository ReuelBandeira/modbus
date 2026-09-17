import 'package:flutter/material.dart';
import 'package:isi_4_0/colors/color.dart';
import 'package:provider/provider.dart';

import '../../../providers/lang.dart';
import '../../../widgets/custom_device_text_field.dart';

class CardTopics extends StatelessWidget {
  final List<String> topicsList;
  final Function(List<String>? value) onChanged;
  final bool isErrored;
  const CardTopics({super.key, required this.topicsList, required this.onChanged,
    required this.isErrored});

  Widget? topicsWidget () {
    List<Widget> list = [];

    for (String topic in topicsList) {
      list.add(
        Padding (
          padding: const EdgeInsets.all(4.0),
          child: SizedBox (
            height: 25,
            child: OutlinedButton.icon(
              onPressed: () {
                topicsList.remove(topic);
                onChanged(topicsList);
              },
              style: OutlinedButton.styleFrom(
                side: const BorderSide(
                    color: CustomColors.neutral500,
                    width: 1
                ),
                shape: const RoundedRectangleBorder(
                    borderRadius: BorderRadius.all(
                      Radius.circular(8)
                    )
                ),
              ),
              icon: const Icon(
                Icons.cancel_rounded,
                size: 12.0,
                color: CustomColors.neutral700,
              ),
              label: Text(topic,
                style: const TextStyle(
                    fontSize: 13,
                    color: CustomColors.neutral800
                ),
              )
            )
          )
        )
      );
    }

    return Wrap(
        children: list
    );
  }

  @override
  Widget build(BuildContext context) {
    TextEditingController controller = TextEditingController();
    return Consumer<LanguageProvider>(
      builder: (context, language, child) {
        final translator = language.getDataLanguage(language.currentLanguage);
        return SizedBox(
          height: 375,
          width: 300,
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
                      child: Text(translator["newDevices"]["titleTopicsCard"],
                          style: const TextStyle(
                              fontSize: 16,
                              color: CustomColors.neutral800,
                              fontWeight: FontWeight.bold))),
                  Container(
                      alignment: Alignment.topLeft,
                      child: Text(translator["newDevices"]["subtitleTopicsCard"],
                          style: const TextStyle(fontSize: 12, color: CustomColors.neutral700)
                      )
                  ),
                  Form(
                    child: Container(
                      height: 290,
                      alignment: Alignment.bottomCenter,
                      padding: const EdgeInsets.only(top: 8.0),
                      child: Column(
                        mainAxisAlignment: MainAxisAlignment.spaceBetween,
                        children: [
                          Container(
                            height: 168,
                            width: 300,
                            color: CustomColors.background700,
                            child: SingleChildScrollView(
                              child: topicsWidget()
                            ),
                          ),
                          Container(
                              margin: const EdgeInsets.only(top: 6.0),
                              child: CustomDeviceTextField(
                                  controller: controller,
                                  labelText: translator["newDevices"]["titleTopicsCard"]
                              )
                          ),
                          Container(
                            height: 40,
                            margin: const EdgeInsets.only(top: 6.0),
                            child: TextButton(
                                onPressed: () {
                                  if (controller.text != '') {
                                    topicsList.add(controller.text);
                                    onChanged(topicsList);
                                    controller.clear();
                                  }
                                },
                                child: Text(translator["newDevices"]["addButton"])),
                          ),
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
