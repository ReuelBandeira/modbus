import 'package:flutter/material.dart';
import 'package:isi_4_0/views/devices/new_devices.dart';

class PageViewExtraData extends StatefulWidget {
  const PageViewExtraData({super.key, this.onTap});
  final Function()? onTap;

  @override
  State<PageViewExtraData> createState() => __PageViewExtraDataState();
}

class __PageViewExtraDataState extends State<PageViewExtraData> {
  PageController pageController = PageController();
  dynamic bufferDatareceive;

  @override
  Widget build(BuildContext context) {
    void receiveDataFromChild(dynamic data) {
      setState(() {
        bufferDatareceive = data;
      });
    }

    return PageView(
      physics: const NeverScrollableScrollPhysics(),
      controller: pageController,
      scrollDirection: Axis.horizontal,
      children: [
        NewDevices(
          onTap: widget.onTap,
          dataCallback: receiveDataFromChild,
          extraData: () {
            // print('Data: $bufferDatareceive');
            pageviewNavigator(pageController, 1);
          },
        )
      ],
    );
  }

  void pageviewNavigator(PageController pageController, int position) {
    pageController.animateToPage(position,
        duration: const Duration(milliseconds: 1000), curve: Curves.easeIn);
  }
}
