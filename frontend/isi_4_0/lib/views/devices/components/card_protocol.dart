import 'package:flutter/material.dart';
import 'package:isi_4_0/colors/color.dart';
import 'package:isi_4_0/utils/data_protocols.dart';
import 'package:isi_4_0/widgets/custom_rounded_button.dart';
import 'package:provider/provider.dart';

import '../../../providers/lang.dart';

class CardProtocol extends StatelessWidget {
  const CardProtocol({super.key, required this.protocols, this.onChanged, required this.value, required this.isErrored});
  final List<DataProtocols> protocols;
  final Function(String? value)? onChanged;
  final String? value;
  final bool isErrored;
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
              child: Padding(
                padding: const EdgeInsets.all(16),
                child: Column(
                  mainAxisAlignment: MainAxisAlignment.start,
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Container(
                        alignment: Alignment.bottomLeft,
                        padding: const EdgeInsets.only(bottom: 8.0),
                        child: Text(translator["newDevices"]["titleProtocolsCard"],
                            style: const TextStyle(
                                fontSize: 16,
                                color: CustomColors.neutral800,
                                fontWeight: FontWeight.bold))),
                    Container(
                        alignment: Alignment.topLeft,
                        child: Text(translator["newDevices"]["subtitleProtocolsCard"],
                            style: const TextStyle(fontSize: 12, color: CustomColors.neutral700))),
                    Container(
                      height: 60,
                      alignment: Alignment.bottomCenter,
                      child: Container(
                        height: 36,
                        width: 250,
                        alignment: Alignment.center,
                        decoration: BoxDecoration(
                          color: CustomColors.background700,
                          border: Border.all(
                            color: CustomColors.neutral500,
                            width: 1.0,
                          ),
                          borderRadius: BorderRadius.circular(4.0),
                        ),
                        child: Padding(
                          padding: const EdgeInsets.only(left: 12.0, right: 12.0),
                          child: Theme(
                            data: Theme.of(context).copyWith(
                              // splashColor: Colors.transparent,
                              // highlightColor: Colors.transparent,
                              // hoverColor: Colors.transparent,
                              // indicatorColor: Colors.transparent,
                            ),
                            child: DropdownButton<String>(
                              hint: Text(
                                translator["newDevices"]["fieldSelect"],
                                style: const TextStyle(fontSize: 16),
                              ),
                              value: value,
                              isExpanded: true,
                              underline: Container(),
                              style: const TextStyle(
                                  color: CustomColors.neutral700,
                                  fontSize: 16
                              ),
                              onChanged: (newValue) {
                                if (onChanged != null) {
                                  onChanged!(newValue);
                                }
                              },
                              items:
                              protocols.map<DropdownMenuItem<String>>((DataProtocols protocol) {
                                return DropdownMenuItem<String>(
                                  value: protocol.protocol,
                                  child: Text(protocol.alias),
                                );
                              }).toList(),
                            ),
                          ),
                        ),
                      ),
                    ),
                  ],
                ),
              ),
            ),
          );
        }
    );
  }
}
