import 'package:flutter/material.dart';
import 'package:isi_4_0/colors/color.dart';

class CustomRowListDevice extends StatelessWidget {
  const CustomRowListDevice(
      {super.key,
      required this.device,
      required this.hint,
      required this.tooltipsMessage,
      this.buttonsName,
      this.onSelected});
  final dynamic device;
  final String hint;
  final String tooltipsMessage;
  final dynamic buttonsName;
  final void Function(String)? onSelected;
  @override
  Widget build(BuildContext context) {
    int id = device['id'];
    List<dynamic> dataTopics = device['topics'];
    final topics = dataTopics.map((topic) => topic.toString()).toList();
    String? selectedValue;

    const double height = 40;
    const double width = 180;
    return Container(
      height: height,
      margin: const EdgeInsets.only(top: 2.0),
      padding: const EdgeInsets.only(left: 4.0),
      child: SingleChildScrollView(
        scrollDirection: Axis.horizontal,
        child: Row(
          children: [
            Container(
              width: 4,
              color: device['status']
                  ? CustomColors.deviceOnline
                  : CustomColors.deviceOffOnline,
            ),
            Container(
              height: height,
              width: width - 4,
              alignment: Alignment.centerLeft,
              padding: const EdgeInsets.only(left: 8.0),
              decoration: const BoxDecoration(
                  border: Border(
                      bottom: BorderSide(
                          color: CustomColors.lineColors, width: 2))),
              child: Text(device['address']),
            ),
            Container(
              height: height,
              width: width,
              alignment: Alignment.centerLeft,
              padding: const EdgeInsets.only(left: 8.0),
              decoration: const BoxDecoration(
                  border: Border(
                      bottom: BorderSide(
                          color: CustomColors.lineColors, width: 2))),
              child: Text(device['port']),
            ),
            Container(
              height: height,
              width: width,
              alignment: Alignment.centerLeft,
              padding: const EdgeInsets.only(left: 8.0),
              decoration: const BoxDecoration(
                  border: Border(
                      bottom: BorderSide(
                          color: CustomColors.lineColors, width: 2))),
              child: Text(device['name']),
            ),
            Container(
              height: height,
              width: width,
              alignment: Alignment.centerLeft,
              padding: const EdgeInsets.only(left: 8.0),
              decoration: const BoxDecoration(
                  border: Border(
                      bottom: BorderSide(
                          color: CustomColors.lineColors, width: 2))),
              child: Text(device['protocol']),
            ),
            Container(
              height: height,
              width: width,
              alignment: Alignment.centerLeft,
              padding: const EdgeInsets.only(left: 8.0),
              decoration: const BoxDecoration(
                  border: Border(
                      bottom: BorderSide(
                          color: CustomColors.lineColors, width: 2))),
              child: Text(device['readingTime'].toString()),
            ),
            Container(
              height: height,
              width: width,
              padding: const EdgeInsets.only(left: 8.0),
              decoration: const BoxDecoration(
                  border: Border(
                      bottom: BorderSide(
                          color: CustomColors.lineColors, width: 2))),
              child: Row(
                children: [
                  SizedBox(
                    height: height,
                    width: width - (height + 8),
                    child: DropdownButton<String>(
                      hint: Text(
                        hint,
                        style: const TextStyle(fontSize: 13),
                      ),
                      value: selectedValue,
                      isExpanded: true,
                      underline: const SizedBox.shrink(),
                      onChanged: (newValue) {},
                      items:
                          topics.map<DropdownMenuItem<String>>((String value) {
                        return DropdownMenuItem<String>(
                          value: value,
                          child: Text(value),
                        );
                      }).toList(),
                    ),
                  ),
                  Align(
                    alignment: Alignment.center,
                    child: SizedBox(
                      height: 35,
                      width: 35,
                      child: PopupMenuButton<String>(
                        tooltip: tooltipsMessage,
                        onSelected: onSelected,
                        itemBuilder: (BuildContext context) =>
                            <PopupMenuEntry<String>>[
                          PopupMenuItem<String>(
                            value: 'edit',
                            child:
                                Text(buttonsName['dictionary']['editDevice']),
                          ),
                          PopupMenuItem<String>(
                            value: 'remove',
                            child:
                                Text(buttonsName['dictionary']['removeDevice']),
                          ),
                        ],
                        child: const Icon(
                          Icons.more_vert,
                          size: 20,
                        ),
                      ),
                    ),
                  )
                ],
              ),
            )
          ],
        ),
      ),
    );
  }
}
