import 'package:flutter/material.dart';
import 'package:isi_4_0/colors/color.dart';
import 'package:isi_4_0/widgets/custom_device_text_field.dart';
import 'package:provider/provider.dart';

import '../../../providers/lang.dart';

class CardTextFieldProtocol extends StatelessWidget {
  final TextEditingController textControllerIp;
  final TextEditingController textControllerPort;
  final TextEditingController textControllerName;
  final TextEditingController textControllerReadTime;
  final bool isErrored;
  const CardTextFieldProtocol({
    super.key,
    required this.textControllerIp,
    required this.textControllerPort,
    required this.textControllerName,
    required this.textControllerReadTime,
    required this.isErrored
  });
  @override
  Widget build(BuildContext context) {
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
                        child: Text(translator["newDevices"]["titleConfigurationsCard"],
                            style: const TextStyle(
                                fontSize: 16,
                                color: CustomColors.neutral800,
                                fontWeight: FontWeight.bold))),
                    Container(
                        alignment: Alignment.topLeft,
                        child: Text(translator["newDevices"]["subtitleConfigurationsCard"],
                            style: const TextStyle(fontSize: 12, color: CustomColors.neutral700))),
                    Form(
                        // child: Container(
                        //   height: 290,
                        //   alignment: Alignment.topCenter,
                        //   padding: const EdgeInsets.only(top: 8.0),
                          child: Expanded(
                            child: SingleChildScrollView(
                              scrollDirection: Axis.vertical,
                              child: Column(
                                mainAxisAlignment: MainAxisAlignment.spaceBetween,
                                children: [
                                  Container(
                                      margin: const EdgeInsets.only(top: 6.0),
                                      child: CustomDeviceTextField(
                                        controller: textControllerName,
                                        labelText: translator["newDevices"]["fieldName"],
                                        isRequired: true,
                                      )
                                  ),
                                  Container(
                                      margin: const EdgeInsets.only(top: 6.0),
                                      child: CustomDeviceTextField(
                                        controller: textControllerIp,
                                        labelText: translator["newDevices"]["fieldIpAddress"],
                                        ipOnly: true,
                                        isRequired: true,
                                      )
                                  ),
                                  Container(
                                      margin: const EdgeInsets.only(top: 6.0),
                                      child: CustomDeviceTextField(
                                        controller: textControllerPort,
                                        labelText: translator["newDevices"]["fieldPort"],
                                        digitsOnly: true,
                                        isRequired: true,
                                      )
                                  ),
                                  Container(
                                      margin: const EdgeInsets.only(top: 6.0),
                                      child: CustomDeviceTextField(
                                        controller: textControllerReadTime,
                                        labelText: translator["newDevices"]["fieldReadingTime"],
                                        digitsOnly: true,
                                        isRequired: true,
                                      )
                                  ),
                                ],
                              ),
                            )
                          )
                        // )
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
