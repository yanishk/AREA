import 'package:flutter/material.dart';
import 'package:auto_size_text/auto_size_text.dart';

import "custom_widget.dart";
import "globalData.dart";

class SelectServicePage extends StatelessWidget {
  const SelectServicePage({
    super.key,
    this.isAnAction = true,
  });
  final bool isAnAction;

  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      theme: ThemeData(scaffoldBackgroundColor: const Color.fromARGB(255, 255, 255, 255)),
      home: APICall(isAnAction: isAnAction),
    );
  }
}

class APICall extends StatefulWidget {
  const APICall({
    super.key,
    this.isAnAction = true,
  });
  final bool isAnAction;

  @override
  SelectService createState() => SelectService(isAnAction: isAnAction);
}

class SelectService extends State<APICall> {
  SelectService({
    required this.isAnAction,
  });
  final bool isAnAction;
  String title = "Select action service";
  List<String> serviceNameList = Globaldata.actionServicesName;
  List<String> serviceColorList = Globaldata.actionServicesColor;
  List<String> serviceIconsList = Globaldata.actionServicesIcons;

  Future<void> setReactionService() async {
    setState(() {
      title = "Select reaction service";
      serviceNameList = Globaldata.reactionServicesName;
      serviceColorList = Globaldata.reactionServicesColor;
      serviceIconsList = Globaldata.reactionServicesIcons;
    });
  }

  @override
  Widget build(BuildContext context) {
    if (isAnAction != true) {
      setReactionService();
    }
    List<Widget> services = [];
    for (int i = 0; i < serviceNameList.length; i++) {
      services.add(
        Row(
          children: [
            MySquareServicesCard(
              servicesIndex: i,
              isAnAction: isAnAction,
              serviceColorList: serviceColorList,
              serviceNameList: serviceNameList,
              serviceIconsList: serviceIconsList,
            ),
            // MySquareServicesCard(
            //   servicesIndex: i + 1,
            //   isAnAction: isAnAction,
            // ),
          ],
        ),
      );
    }

    return Scaffold(
      body: SingleChildScrollView(
        child: Column(
          children: [
            Padding(padding: EdgeInsets.only(top: 30)),
            SelectServiceTopBar(title: title),

            MyDivider(),

            Column(
              children: services,
            ),
          ],
        ),
      ),
    );
  }
}

class SelectServiceTopBar extends StatelessWidget {
  const SelectServiceTopBar({
    super.key,
    required this.title,
  });

  final String title;

  @override
  Widget build(BuildContext context) {
    return Stack(
      alignment: Alignment.center,
      children: [
        Positioned(
          left: 0,
          child: IconButton(
            onPressed: () {Navigator.pop(Globaldata.myContext);},
            icon: Icon(Icons.arrow_back, size: 32),
          ),
        ),
        Center(
          child: AutoSizeText(title, style: TextStyle( fontSize: 30, fontWeight: FontWeight.bold, color: Colors.black,), maxLines: 1,),
        ),
      ],
    );
  }
}



class MySquareServicesCard extends StatelessWidget {
  const MySquareServicesCard({
    super.key,
    required this.servicesIndex,
    this.isAnAction = true,
    required this.serviceNameList,
    required this.serviceColorList,
    required this.serviceIconsList,
  });

  final int servicesIndex;
  final bool isAnAction;
  final List<String> serviceNameList;
  final List<String> serviceColorList;
  final List<String> serviceIconsList;

  @override
  Widget build(BuildContext context) {
    if (servicesIndex >= Globaldata.actionServicesName.length) {
      return MyDivider();
    }
    return Padding(
      padding: const EdgeInsets.symmetric(horizontal: 15, vertical: 10),
      child: LayoutBuilder(
        builder: (context, constraints) {
          return Container(
            decoration: BoxDecoration(
              color: HexColor(serviceColorList[servicesIndex]),
              border: Border.all(
                color: HexColor(serviceColorList[servicesIndex]),
                width: 2.0,
              ),
              borderRadius: BorderRadius.circular(4.0),
            ),

            child: InkWell(
              onTap: () {
                if (isAnAction) {
                  Globaldata.firstService = servicesIndex;
                } else {
                  Globaldata.secondService = servicesIndex;
                }
                Navigator.pop(Globaldata.myContext, serviceNameList[servicesIndex]);
              },
              child: Column(
                children: [
                  SizedBox(
                    width: MediaQuery.sizeOf(Globaldata.myContext).width / 2.5,
                  ),
                  Padding(
                    padding: const EdgeInsets.symmetric(horizontal: 5, vertical: 30),
                    child: Column(
                      children: [
                        FadeInImage(
                          image: NetworkImage("https://mathisbrehm.fr/${serviceIconsList[servicesIndex]}"),
                          placeholder: AssetImage('assets/icon/noImage.png'),
                          imageErrorBuilder: (context, error, stackTrace) {
                            return Image.asset(
                              'assets/icon/noImage.png',
                              width: 75,
                              height: 75,
                            );
                          },
                          width: 75,
                          height: 75,
                        ),
                        AutoSizeText(serviceNameList[servicesIndex], style: const TextStyle(fontSize: 25, fontWeight: FontWeight.bold, color: Colors.white), maxLines: 1,)
                      ]
                    ),
                  ),
                ],
              ),
            ),
          );
        },
      ),
    );
  }
}

class HexColor extends Color {
  static int _getColorFromHex(String hexColor) {
    hexColor = hexColor.toUpperCase().replaceAll("#", "");
    if (hexColor.length == 6) {
      hexColor = "FF$hexColor";
    }
    return int.parse(hexColor, radix: 16);
  }

  HexColor(final String hexColor) : super(_getColorFromHex(hexColor));
}