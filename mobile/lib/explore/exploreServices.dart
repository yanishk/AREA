import 'package:flutter/material.dart';

import "../custom_widget.dart";
import "../footer_page.dart";
import "../globalData.dart";
import "topBarExplore.dart";
import "exploreWidget.dart";

class ExplorePageServices extends StatelessWidget {
  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      theme: new ThemeData(
          scaffoldBackgroundColor: const Color.fromARGB(255, 255, 255, 255)),
      home: ExplorePageServicesStateful(),
    );
  }
}

class ExplorePageServicesStateful extends StatefulWidget {
  @override
  ExplorePageServicesMain createState() => ExplorePageServicesMain();
}

class ExplorePageServicesMain extends State<ExplorePageServicesStateful> {
  int serviceIndex = 0;
  List<String> serviceNameList = Globaldata.actionServicesName + Globaldata.reactionServicesName;
  List<String> serviceColorList = Globaldata.actionServicesColor + Globaldata.reactionServicesColor;
  List<String> serviceIconsList = Globaldata.actionServicesIcons + Globaldata.reactionServicesIcons;

  Widget serviceTiles(int index) {
    return ServicesCardParam(
      name: serviceNameList[index],
      color: HexColor(serviceColorList[index]),
      imageUrl: serviceIconsList[index],
    );
  }

  @override
  Widget build(BuildContext context) {
    List<Widget> services = [];
    for (serviceIndex = 0;
        serviceIndex < serviceNameList.length;
        serviceIndex++) {
      services.add(serviceTiles(serviceIndex));
    }
    return Scaffold(
      body: SingleChildScrollView(
        child: Column(
          children: [
            Padding(padding: EdgeInsets.only(top: 30)),
            TopBarExplore(
              currentIndex: 2,
            ),
            MyDivider(),
            MySearchBar(),
            Column(
              children: services,
            ),
          ],
        ),
      ),
      bottomNavigationBar: FooterNavigationBar(currentIndex: 1),
    );
  }
}